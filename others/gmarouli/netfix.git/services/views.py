from django.shortcuts import render, redirect, get_object_or_404
from django.http import HttpResponseRedirect
from django.contrib import messages
from django.contrib.auth.decorators import login_required
from django.db.models import Count
from django.views.decorators.http import require_http_methods, require_POST
from django.db import transaction
from django.core.exceptions import ValidationError
import logging
from django.http import Http404

from users.models import Company, Customer, User

from .models import Service, ServiceRequest
from .services import CreateNewService
from .forms import RequestServiceForm

logger = logging.getLogger(__name__)

def list(request):
    """Display all services"""
    services = Service.objects.all().order_by('-date')
    return render(request, 'services/list.html', {'services': services})


def index(request, id):
    service = Service.objects.get(id=id)
    return render(request, 'services/single_service.html', {'service': service})


@require_http_methods(["GET", "POST"])
def create(request):
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to create services.')
        return redirect('login_user')
    
    # Check if user is a company
    if not request.user.is_company:
        messages.error(request, 'Only companies can create services.')
        logger.warning(f'Non-company user {request.user.username} attempted to create service')
        return redirect('services_list')
    
    # Get the company instance
    try:
        company = Company.objects.get(user=request.user)
    except Company.DoesNotExist:
        messages.error(request, 'Company profile not found.')
        logger.error(f'Company profile not found for user {request.user.username}')
        return redirect('services_list')
    
    # Determine available service field choices based on company type
    if company.field == 'All in One':
        # All in One companies can create services in any category except "All in One"
        choices = [choice for choice in Service.choices if choice[0] != 'All in One']
    else:
        # Other companies can only create services in their specific category
        choices = [(company.field, company.field)]
    
    if request.method == 'POST':
        form = CreateNewService(request.POST, choices=choices)
        if form.is_valid():
            # Additional validation for company field restrictions
            selected_field = form.cleaned_data['field']
            if company.field != 'All in One' and selected_field != company.field:
                messages.error(request, f'You can only create services in your company category: {company.field}')
                logger.warning(f'Company {company.user.username} attempted to create service outside their category')
            elif company.field == 'All in One' and selected_field == 'All in One':
                messages.error(request, 'All in One companies cannot create services in the "All in One" category.')
                logger.warning(f'All in One company {company.user.username} attempted to create service in "All in One" category')
            else:
                try:
                    with transaction.atomic():
                        # Create the service
                        service = Service.objects.create(
                            company=company,
                            name=form.cleaned_data['name'],
                            description=form.cleaned_data['description'],
                            price_hour=form.cleaned_data['price_hour'],
                            field=form.cleaned_data['field']
                        )
                        messages.success(request, f'Service "{service.name}" created successfully!')
                        logger.info(f'Service {service.name} created by company {company.user.username}')
                        return redirect('index', id=service.id)
                except ValidationError as e:
                    messages.error(request, f'Error creating service: {e}')
                    logger.error(f'Service creation failed for {company.user.username}: {e}')
                except Exception as e:
                    messages.error(request, 'An unexpected error occurred. Please try again.')
                    logger.error(f'Unexpected error in service creation: {e}')
        else:
            messages.error(request, 'Please correct the errors below.')
    else:
        form = CreateNewService(choices=choices)
    
    return render(request, 'services/create.html', {'form': form, 'company': company})


def service_field(request, field):
    # Define valid field choices (should match the model choices, excluding "all-in-one")
    valid_fields = [
        'air-conditioner', 'carpentry', 'electricity', 
        'gardening', 'home-machines', 'house-keeping', 'interior-design',
        'locks', 'painting', 'plumbing', 'water-heaters'
    ]
    
    # Check if the field is valid
    if field.lower() not in valid_fields:
        # Render our beautiful custom 404 page instead of raising Http404
        from django.http import HttpResponseNotFound
        return HttpResponseNotFound(render(request, 'main/404.html'))
    
    # Convert URL slug to proper field name for database lookup
    field_name = field.replace('-', ' ').title()
    services = Service.objects.filter(field=field_name)
    
    return render(request, 'services/field.html', {
        'services': services, 
        'field': field_name
    })


@require_http_methods(["GET", "POST"])
def request_service(request, id):
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to request services.')
        return redirect('login_user')
    
    # Get the service or return 404
    try:
        service = get_object_or_404(Service, id=id)
    except ValueError:
        messages.error(request, 'Invalid service ID.')
        return redirect('services_list')
    
    # Check if user is a customer
    if not request.user.is_customer:
        messages.error(request, 'Only customers can request services.')
        logger.warning(f'Non-customer user {request.user.username} attempted to request service {id}')
        return redirect('index', id=service.id)
    
    # Get the customer instance
    try:
        customer = Customer.objects.get(user=request.user)
    except Customer.DoesNotExist:
        messages.error(request, 'Customer profile not found.')
        logger.error(f'Customer profile not found for user {request.user.username}')
        return redirect('index', id=service.id)
    
    if request.method == 'POST':
        form = RequestServiceForm(request.POST)
        if form.is_valid():
            try:
                from datetime import datetime
                
                with transaction.atomic():
                    # Convert start_time string to time object
                    start_time_str = form.cleaned_data['start_time']
                    start_time = datetime.strptime(start_time_str, '%H:%M').time()
                    
                    # Get duration as integer
                    duration_hours = int(form.cleaned_data['duration_hours'])
                    
                    # Create the service request
                    service_request = ServiceRequest.objects.create(
                        customer=customer,
                        service=service,
                        address=form.cleaned_data['address'],
                        service_date=form.cleaned_data['service_date'],
                        start_time=start_time,
                        duration_hours=duration_hours
                    )
                    
                    total_price = duration_hours * service.price_hour
                    end_time = service_request.get_end_time()
                    
                    messages.success(request, f'Service request submitted successfully! Scheduled from {start_time.strftime("%H:%M")} to {end_time.strftime("%H:%M")} ({duration_hours} hour(s)). Total cost: €{total_price}')
                    logger.info(f'Service request created by {customer.user.username} for service {service.name}')
                    return redirect('customer_profile', name=request.user.username)
            except ValueError as e:
                messages.error(request, 'Invalid time format.')
                logger.error(f'Time parsing error in service request: {e}')
            except ValidationError as e:
                messages.error(request, f'Validation error: {e}')
                logger.error(f'Service request validation failed: {e}')
            except Exception as e:
                messages.error(request, 'An unexpected error occurred. Please try again.')
                logger.error(f'Unexpected error in service request: {e}')
        else:
            messages.error(request, 'Please correct the errors below.')
    else:
        form = RequestServiceForm()
    
    return render(request, 'services/request_service.html', {
        'form': form,
        'service': service
    })


@require_POST
def manage_request(request, request_id, action):
    """Accept or reject a service request"""
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to manage service requests.')
        return redirect('login_user')
    
    # Check if user is a company
    if not request.user.is_company:
        messages.error(request, 'Only companies can manage service requests.')
        logger.warning(f'Non-company user {request.user.username} attempted to manage request {request_id}')
        return redirect('services_list')
    
    # Validate action parameter
    if action not in ['accept', 'reject']:
        messages.error(request, 'Invalid action.')
        logger.warning(f'Invalid action {action} attempted by {request.user.username}')
        return redirect('company_profile', name=request.user.username)
    
    # Get the company instance
    try:
        company = Company.objects.get(user=request.user)
    except Company.DoesNotExist:
        messages.error(request, 'Company profile not found.')
        logger.error(f'Company profile not found for user {request.user.username}')
        return redirect('services_list')
    
    # Get the service request
    try:
        service_request = get_object_or_404(ServiceRequest, id=request_id)
    except ValueError:
        messages.error(request, 'Invalid request ID.')
        return redirect('company_profile', name=request.user.username)
    
    # Verify that this request belongs to one of the company's services
    if service_request.service.company != company:
        messages.error(request, 'You are not authorized to manage this request.')
        logger.warning(f'Unauthorized request management attempt by {request.user.username} for request {request_id}')
        return redirect('company_profile', name=request.user.username)
    
    # Check if request is still pending
    if service_request.status != 'pending':
        messages.error(request, 'This request has already been processed.')
        return redirect('company_profile', name=request.user.username)
    
    try:
        with transaction.atomic():
            # Update the status based on action
            if action == 'accept':
                service_request.status = 'accepted'
                service_request.save()
                messages.success(request, f'Service request from {service_request.customer.user.username} has been accepted.')
                logger.info(f'Request {request_id} accepted by {company.user.username}')
            elif action == 'reject':
                service_request.status = 'rejected'
                service_request.save()
                messages.success(request, f'Service request from {service_request.customer.user.username} has been rejected.')
                logger.info(f'Request {request_id} rejected by {company.user.username}')
    except Exception as e:
        messages.error(request, 'An error occurred while processing the request.')
        logger.error(f'Error managing request {request_id}: {e}')
    
    return redirect('company_profile', name=request.user.username)


def most_requested(request):
    """Display most requested services based on request count"""
    from django.db.models import Count
    
    services = Service.objects.annotate(
        request_count=Count('servicerequest')
    ).filter(request_count__gt=0).order_by('-request_count')[:10]
    
    return render(request, 'services/most_requested.html', {'services': services})

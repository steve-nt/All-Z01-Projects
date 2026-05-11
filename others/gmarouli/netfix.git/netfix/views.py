from django.shortcuts import render, get_object_or_404, redirect
from django.contrib.auth.decorators import login_required
from django.contrib import messages
from django.core.exceptions import PermissionDenied
from datetime import date
import logging

from users.models import User, Company, Customer
from services.models import Service, ServiceRequest

logger = logging.getLogger(__name__)


def home(request):
    return render(request, 'users/home.html', {'user': request.user})


def profile_redirect(request, user_id):
    """Redirect from user ID to proper profile URL"""
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to view profiles.')
        return redirect('login_user')
    
    user = get_object_or_404(User, id=user_id)
    
    # Users can only view their own profile
    if request.user != user:
        messages.error(request, 'You can only view your own profile.')
        logger.warning(f'User {request.user.username} attempted to access profile of user ID {user_id}')
        return redirect('customer_profile' if request.user.is_customer else 'company_profile', name=request.user.username)
    
    if user.is_customer:
        return redirect('customer_profile', name=user.username)
    elif user.is_company:
        return redirect('company_profile', name=user.username)
    else:
        # Fallback to customer profile
        return redirect('customer_profile', name=user.username)


def customer_profile(request, name):
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to view your profile.')
        return redirect('login_user')
    
    # Get the user by username
    user = get_object_or_404(User, username=name)
    
    # Security check: Users can only view their own profile
    if request.user != user:
        messages.error(request, 'You can only view your own profile.')
        logger.warning(f'User {request.user.username} attempted to access customer profile of {name}')
        return redirect('customer_profile', name=request.user.username)
    
    # Ensure the requested profile is actually a customer
    if not user.is_customer:
        messages.error(request, 'Invalid profile type.')
        return redirect('/')
    
    # Calculate user age if customer
    user_age = None
    service_requests = []
    
    try:
        customer = Customer.objects.get(user=user)
        # Calculate age
        today = date.today()
        user_age = today.year - customer.birth.year - ((today.month, today.day) < (customer.birth.month, customer.birth.day))
        
        # Get all service requests for this customer
        service_requests = ServiceRequest.objects.filter(customer=customer).order_by('-request_date')
        
    except Customer.DoesNotExist:
        messages.error(request, 'Customer profile not found.')
        logger.error(f'Customer profile not found for user {user.username}')
        return redirect('/')
    
    context = {
        'user': user,
        'user_age': user_age,
        'sh': service_requests,  # Using 'sh' to match the template variable
    }
    
    return render(request, 'users/profile.html', context)


def company_profile(request, name):
    # Check if user is authenticated
    if not request.user.is_authenticated:
        messages.info(request, 'Please log in to view profiles.')
        return redirect('login_user')
    
    # Get the user by username
    user = get_object_or_404(User, username=name)
    
    # Security check: Users can only view their own profile
    if request.user != user:
        messages.error(request, 'You can only view your own profile.')
        logger.warning(f'User {request.user.username} attempted to access company profile of {name}')
        return redirect('customer_profile' if request.user.is_customer else 'company_profile', name=request.user.username)
    
    # Ensure the requested profile is actually a company
    if not user.is_company:
        messages.error(request, 'Invalid profile type.')
        return redirect('/')
    
    try:
        company = Company.objects.get(user=user)
        # Fetch the company's services
        services = Service.objects.filter(company=company).order_by("-date")
        
        # Fetch service requests for this company's services
        service_requests = ServiceRequest.objects.filter(service__company=company).order_by('-request_date')
        
        # Calculate stats
        pending_count = service_requests.filter(status='pending').count()
        accepted_count = service_requests.filter(status='accepted').count()
        rejected_count = service_requests.filter(status='rejected').count()
        
    except Company.DoesNotExist:
        messages.error(request, 'Company profile not found.')
        logger.error(f'Company profile not found for user {user.username}')
        return redirect('/')
    
    # Check if the current user is viewing their own profile (for edit permissions)
    is_own_profile = request.user == user
    
    context = {
        'user': user,
        'services': services,
        'company': company,
        'is_own_profile': is_own_profile,
        'service_requests': service_requests,
        'pending_count': pending_count,
        'accepted_count': accepted_count,
        'rejected_count': rejected_count,
    }

    return render(request, 'users/profile.html', context)


def public_company_view(request, name):
    """Public view of a company's services - accessible to anyone"""
    # Get the user by username
    user = get_object_or_404(User, username=name)
    
    # Ensure the requested profile is actually a company
    if not user.is_company:
        messages.error(request, 'Invalid company profile.')
        return redirect('/')
    
    try:
        company = Company.objects.get(user=user)
        # Fetch the company's services
        services = Service.objects.filter(company=company).order_by("-date")
        
    except Company.DoesNotExist:
        messages.error(request, 'Company profile not found.')
        return redirect('/')
    
    context = {
        'company': company,
        'services': services,
        'is_public_view': True,  # Flag to indicate this is a public view
    }

    return render(request, 'services/public_company.html', context)

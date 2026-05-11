from django.shortcuts import render, redirect
from django.contrib.auth import login, authenticate
from django.views.generic import CreateView, TemplateView
from django.contrib import messages

from .forms import CustomerSignUpForm, CompanySignUpForm, UserLoginForm
from .models import User, Company, Customer


def register(request):
    return render(request, 'users/register.html')


class CustomerSignUpView(CreateView):
    model = User
    form_class = CustomerSignUpForm
    template_name = 'users/register_customer.html'

    def get_context_data(self, **kwargs):
        kwargs['user_type'] = 'customer'
        return super().get_context_data(**kwargs)

    def form_valid(self, form):
        user = form.save()
        login(self.request, user)
        # Redirect customers to services list page
        return redirect('services_list')


class CompanySignUpView(CreateView):
    model = User
    form_class = CompanySignUpForm
    template_name = 'users/register_company.html'

    def get_context_data(self, **kwargs):
        kwargs['user_type'] = 'company'
        return super().get_context_data(**kwargs)

    def form_valid(self, form):
        user = form.save()
        login(self.request, user)
        # Redirect to company profile page after successful registration
        return redirect('company_profile', name=user.username)


def LoginUserView(request):
    if request.method == 'POST':
        form = UserLoginForm(request.POST)
        if form.is_valid():
            email = form.cleaned_data['email']
            password = form.cleaned_data['password']
            
            # Find user by email (since we login with email, not username)
            try:
                user = User.objects.get(email=email)
                username = user.username
            except User.DoesNotExist:
                messages.error(request, 'Invalid email or password.')
                return render(request, 'users/login.html', {'form': form})
            
            # Authenticate using username and password
            user = authenticate(request, username=username, password=password)
            if user is not None:
                login(request, user)
                # Redirect based on user type
                if user.is_customer:
                    return redirect('services_list')  # Customers go to services list
                elif user.is_company:
                    return redirect('company_profile', name=user.username)  # Companies go to their profile
                else:
                    return redirect('/')  # Fallback to homepage
            else:
                messages.error(request, 'Invalid email or password.')
    else:
        form = UserLoginForm()
    
    return render(request, 'users/login.html', {'form': form})

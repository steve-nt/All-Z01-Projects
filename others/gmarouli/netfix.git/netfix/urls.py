"""netfix URL Configuration

The `urlpatterns` list routes URLs to views. For more information please see:
    https://docs.djangoproject.com/en/1.11/topics/http/urls/
Examples:
Function views
    1. Add an import:  from my_app import views
    2. Add a URL to urlpatterns:  url(r'^$', views.home, name='home')
Class-based views
    1. Add an import:  from other_app.views import Home
    2. Add a URL to urlpatterns:  url(r'^$', Home.as_view(), name='home')
Including another URLconf
    1. Import the include() function: from django.conf.urls import url, include
    2. Add a URL to urlpatterns:  url(r'^blog/', include('blog.urls'))
"""
from django.contrib import admin
from django.urls import include, path, re_path
from django.conf import settings
from netfix import views as v
from main.views import custom_404_view

urlpatterns = [
    path('admin/', admin.site.urls),
    path('users/', include('users.urls')),
    path('services/', include('services.urls')),
    path('', include('main.urls')),
    path('<int:user_id>/', v.profile_redirect, name='profile_redirect'),
    path('customer/<slug:name>', v.customer_profile, name='customer_profile'),
    path('company/<slug:name>', v.company_profile, name='company_profile'),
    path('company-public/<slug:name>', v.public_company_view, name='public_company'),
    # Catch-all pattern for 404 - must be last
    re_path(r'^.*/$', custom_404_view, name='catch_all_404'),
]

# Custom error handlers
handler404 = 'main.views.custom_404'
handler500 = 'main.views.custom_500'
handler403 = 'main.views.custom_403'

# Development error page testing (only in DEBUG mode)
if settings.DEBUG:
    from django.views.generic import TemplateView
    urlpatterns += [
        path('test/404/', TemplateView.as_view(template_name='main/404.html'), name='test_404'),
        path('test/500/', TemplateView.as_view(template_name='main/500.html'), name='test_500'),
        path('test/403/', TemplateView.as_view(template_name='main/403.html'), name='test_403'),
    ]

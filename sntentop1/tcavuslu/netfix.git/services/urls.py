from django.urls import path
from . import views as v

urlpatterns = [
    path('', v.list, name='services_list'),
    path('create/', v.create, name='services_create'),
    path('company/dashboard/', v.company_dashboard, name='company_dashboard'),
    path('request/<int:request_id>/<str:action>/', v.manage_request, name='manage_request'),
    path('most-requested/', v.most_requested, name='most_requested_services'),
    path('<int:id>', v.index, name='index'),
    path('<int:id>/request_service/', v.request_service, name='request_service'),
    path('<slug:field>/', v.service_field, name='services_field'),
]

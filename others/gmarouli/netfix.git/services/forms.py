from django import forms
from datetime import date

from users.models import Company


class RequestServiceForm(forms.Form):
    # Start time choices (9 AM to 6 PM)
    START_TIME_CHOICES = [
        ('09:00', '9:00 AM'),
        ('10:00', '10:00 AM'),
        ('11:00', '11:00 AM'),
        ('12:00', '12:00 PM'),
        ('13:00', '1:00 PM'),
        ('14:00', '2:00 PM'),
        ('15:00', '3:00 PM'),
        ('16:00', '4:00 PM'),
        ('17:00', '5:00 PM'),
        ('18:00', '6:00 PM'),
    ]
    
    # Duration choices (1-4 hours)
    DURATION_CHOICES = [
        (1, '1 hour'),
        (2, '2 hours'),
        (3, '3 hours'),
        (4, '4 hours'),
    ]
    
    address = forms.CharField(
        max_length=200, 
        required=True,
        widget=forms.TextInput(attrs={
            'placeholder': 'Enter service address',
            'class': 'form-input'
        })
    )
    
    service_date = forms.DateField(
        widget=forms.DateInput(attrs={
            'type': 'date',
            'class': 'form-input'
        }),
        required=True,
        label='Service Date'
    )
    
    start_time = forms.ChoiceField(
        choices=START_TIME_CHOICES,
        widget=forms.Select(attrs={'class': 'form-input'}),
        required=True,
        label='Start Time'
    )
    
    duration_hours = forms.ChoiceField(
        choices=DURATION_CHOICES,
        widget=forms.Select(attrs={'class': 'form-input'}),
        required=True,
        label='Duration'
    )

    def clean_service_date(self):
        service_date = self.cleaned_data['service_date']
        today = date.today()
        
        if service_date < today:
            raise forms.ValidationError("Service date cannot be in the past. Please select today or a future date.")
        
        return service_date

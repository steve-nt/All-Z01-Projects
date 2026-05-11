# DEPRECATED: This file is no longer used.
# All forms have been moved to forms.py
# This file can be safely deleted.

# The CreateNewService form is now properly defined in:
# services/forms.py (with Bootstrap styling and proper widgets)

from django import forms

from users.models import Company


class CreateNewService(forms.Form):
    name = forms.CharField(
        max_length=40,
        widget=forms.TextInput(attrs={'class': 'form-control'})
    )
    description = forms.CharField(
        widget=forms.Textarea(attrs={'class': 'form-control', 'rows': 4}), 
        label='Description'
    )
    price_hour = forms.DecimalField(
        decimal_places=2, 
        max_digits=5, 
        min_value=0.00,
        widget=forms.NumberInput(attrs={'class': 'form-control', 'step': '0.01'})
    )
    field = forms.ChoiceField(
        required=True,
        widget=forms.Select(attrs={'class': 'form-select'})
    )

    def __init__(self, *args, choices='', ** kwargs):
        super(CreateNewService, self).__init__(*args, **kwargs)
        # adding choices to fields
        if choices:
            self.fields['field'].choices = choices
        # adding placeholders to form fields
        self.fields['name'].widget.attrs['placeholder'] = 'Enter Service Name'
        self.fields['description'].widget.attrs['placeholder'] = 'Enter Description'
        self.fields['price_hour'].widget.attrs['placeholder'] = 'Enter Price per Hour (€)'

        self.fields['name'].widget.attrs['autocomplete'] = 'off'


class RequestServiceForm(forms.Form):
    pass

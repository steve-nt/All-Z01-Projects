from django.shortcuts import render
from django.contrib.auth import logout as django_logout
from django.http import HttpResponseNotFound, HttpResponseServerError, HttpResponseForbidden


def home(request):
    return render(request, "main/home.html", {})


def logout(request):
    django_logout(request)
    return render(request, "main/logout.html")


def custom_404(request, exception):
    """Custom 404 error page"""
    return HttpResponseNotFound(render(request, 'main/404.html').content)


def custom_500(request):
    """Custom 500 error page"""
    return HttpResponseServerError(render(request, 'main/500.html').content)


def custom_403(request, exception=None):
    """Custom 403 Forbidden error page"""
    return HttpResponseForbidden(render(request, 'main/403.html', status=403))


def test_403(request):
    """Test function to trigger 403 error"""
    return HttpResponseForbidden(render(request, 'main/403.html'))


def custom_404_view(request):
    """Custom view to handle unmatched URLs and show our beautiful 404 page"""
    return HttpResponseNotFound(render(request, 'main/404.html'))

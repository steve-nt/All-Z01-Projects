from django.conf import settings
from django.http import Http404, HttpResponseNotFound, HttpResponseServerError, HttpResponseForbidden
from django.shortcuts import render
from django.utils.deprecation import MiddlewareMixin


class CustomErrorPageMiddleware(MiddlewareMixin):
    """
    Middleware to show custom error pages even in DEBUG mode.
    This is for development testing only - remove in production.
    """
    
    def process_exception(self, request, exception):
        """Handle exceptions and show custom error pages if enabled"""
        
        # Only apply custom error pages if we're in DEBUG mode and the URL suggests we want to test them
        if not settings.DEBUG:
            return None
            
        # Check if this is a 404 error and we want custom handling
        if isinstance(exception, Http404):
            # Show custom 404 page for specific paths or if requested
            if (request.path.startswith('/services/') or 
                request.GET.get('show_custom_404') == 'true'):
                return HttpResponseNotFound(
                    render(request, 'main/404.html').content,
                    content_type='text/html'
                )
        
        # Let Django handle other exceptions normally in DEBUG mode
        return None 
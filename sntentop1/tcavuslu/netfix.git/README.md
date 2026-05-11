NetFix - Home Services Marketplace (Django)
=========================================

Overview
--------
NetFix is a Django-based web app that connects customers with service companies (e.g., plumbing, electricity, carpentry). Customers can browse services and request appointments; companies can create services and manage incoming requests via a dashboard.

Objectives
----------
- Provide simple onboarding for two roles: Customer and Company
- Let companies publish services with category, description, and hourly price
- Let customers request a service for a given date/time and duration
- Track service requests with status: pending, accepted, rejected
- Offer a company dashboard to accept or reject requests
- Display most requested services

Tech Stack
---------
- Python 3.x, Django 1.11.x
- SQLite (default) for local development
- Django templates, messages framework

Project Structure
-----------------
- `netfix/`: project settings and root URL routing
- `main/`: homepage and custom error views/templates
- `users/`: custom `User`, `Customer`, `Company` models, registration and login
- `services/`: `Service` and `ServiceRequest` models, listing, creation, requests, dashboard
- `static/`: static assets (CSS/images)
- `templates/`: app-specific templates under `main/`, `users/`, `services/`

Core Models
-----------
- `users.User`: extends Django user with `is_company`, `is_customer`, and unique `email`
- `users.Customer`: one-to-one with `User`, stores `birth`
- `users.Company`: one-to-one with `User`, stores `field` and `rating`
- `services.Service`: belongs to `Company`; fields include `name`, `description`, `price_hour`, `field`, `rating`
- `services.ServiceRequest`: belongs to `Customer` and `Service`; includes address, date/time, duration, computed `price`, and `status`

Key Features
------------
- Registration flows for Customers and Companies
- Email-based login (maps to username under the hood)
- Service creation with field/category restrictions by company type
- Customer service requests with validation and total price calculation
- Company dashboard to accept/reject requests
- Most requested services view

Routes (high-level)
-------------------
- `/` → main homepage
- `/users/` → registration and login
  - `users:`
    - `''` → register landing
    - `company/` → company registration
    - `customer/` → customer registration
    - `login/` → login
- `/services/` → service browsing and actions
  - `''` → services list
  - `create/` → create a new service (company only)
  - `company/dashboard/` → manage requests (company only)
  - `most-requested/` → top requested services
  - `<int:id>` → single service page
  - `<int:id>/request_service/` → request a service (customer only)
  - `<slug:field>/` → list services by field/category (slug)
- Profiles
  - `/customer/<slug:name>` → customer profile (self-only)
  - `/company/<slug:name>` → company profile (self-only)

Custom Error Pages
------------------
Custom 403, 404, 500 pages are provided under `main/templates/main/` and wired in `netfix/urls.py`.

Local Development Setup
-----------------------
1) Prerequisites
- Python 3.x
- pip, virtualenv (recommended)

2) Clone and create a virtual environment
```bash
git clone <your-fork-or-repo-url>
cd netfix
python3 -m venv .venv
source .venv/bin/activate
```

3) Install dependencies
```bash
pip install -r requirements.txt
```

If `requirements.txt` is missing, install minimal versions:
```bash
pip install "Django==1.11.29"
```

4) Database migrations
```bash
python manage.py migrate
```

5) Create a superuser (optional, for admin)
```bash
python manage.py createsuperuser
```

6) Run the server
```bash
python manage.py runserver
```

7) Access the app
- App: `http://127.0.0.1:8000/`
- Admin: `http://127.0.0.1:8000/admin/`

Seeding Data (optional)
-----------------------
- Register a Company and a Customer from the UI:
  - `http://127.0.0.1:8000/users/` → choose Company or Customer
- As a company, create a service at `/services/create/`
- As a customer, request a service at `/services/<id>/request_service/`

Role Flows
----------
- Customer
  - Register and log in
  - Browse `/services/`, open a service, request a date/time/duration
  - View own profile at `/customer/<your-username>` (includes requests list)
- Company
  - Register and log in
  - Create services at `/services/create/` (field restrictions apply)
  - Manage incoming requests at `/services/company/dashboard/`

Environment & Settings
----------------------
- `DEBUG=True` in local settings; change to `False` for production
- Logging writes to `django.log` at project root
- Static files are served from `/static/` during development; see `STATICFILES_DIRS`
- Database: SQLite (`db.sqlite3`) by default

Testing
-------
Run unit tests (if any are added):
```bash
python manage.py test
```

Troubleshooting
---------------
- Login uses email field but authenticates via username internally; ensure the email exists
- Ensure you’re logged in with the correct role for restricted views
- For date/time validation errors when requesting services, verify formats and future dates

Security Notes
--------------
- Keep `SECRET_KEY` out of source control for production
- Set `DEBUG=False` and configure `ALLOWED_HOSTS` in production
- Use HTTPS and set `SESSION_COOKIE_SECURE`/`CSRF_COOKIE_SECURE=True` in production


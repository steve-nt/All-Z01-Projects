# 🏠 Netfix - Home Services Platform

[![Django](https://img.shields.io/badge/Django-5.2.3-green.svg)](https://www.djangoproject.com/)
[![Python](https://img.shields.io/badge/Python-3.8+-blue.svg)](https://www.python.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![Status](https://img.shields.io/badge/Status-Production%20Ready-brightgreen.svg)]()

A modern, full-featured Django web application that connects customers with verified service companies for home maintenance and improvement services. Built with a focus on user experience, security, and scalability.


## 🎯 Overview

Netfix is a comprehensive home services marketplace that bridges the gap between customers seeking reliable home services and professional service companies. The platform provides a seamless experience for both parties with advanced features like real-time scheduling, automated pricing, and comprehensive request management.

### Key Benefits

- **For Customers**: Easy service discovery, transparent pricing, and convenient scheduling
- **For Companies**: Business growth opportunities, streamlined request management, and professional profile showcasing
- **For Platform**: Scalable architecture, secure user management, and comprehensive analytics

## ✨ Features

### 🏠 Customer Features

#### Service Discovery & Browsing
- **Comprehensive Service Catalog**: Browse services across 12 different categories
- **Advanced Filtering**: Filter by service type, price range, and company rating
- **Service Details**: View detailed service information including pricing, descriptions, and company profiles
- **Most Requested Services**: Discover trending and popular services in your area

#### Service Request Management
- **Intelligent Scheduling**: Book services with specific dates, times, and durations
- **Automated Pricing**: Real-time price calculation based on hourly rates and duration
- **Request Tracking**: Monitor request status (Pending, Accepted, Rejected) in real-time
- **Request History**: Complete history of all service requests with detailed information

#### User Experience
- **Personalized Profiles**: Custom customer profiles with request history and preferences
- **Responsive Design**: Optimized for desktop, tablet, and mobile devices
- **Real-time Notifications**: Instant updates on request status changes
- **Secure Authentication**: Email-based login with robust security measures

### 🏢 Company Features

#### Service Management
- **Service Creation**: Create and manage detailed service listings with pricing
- **Category Management**: Offer services in specific categories based on company specialization
- **Service Updates**: Modify service details, pricing, and availability
- **Service Analytics**: Track service performance and customer engagement

#### Request Management
- **Request Dashboard**: Comprehensive view of all incoming service requests
- **Status Management**: Accept or reject requests with detailed feedback
- **Request Analytics**: Track request patterns and business metrics
- **Customer Communication**: Direct communication channel with customers

#### Business Tools
- **Company Profiles**: Professional company profiles with service showcase
- **Public Pages**: Public-facing company pages for customer discovery
- **Rating System**: Customer rating and review system
- **Business Insights**: Analytics and reporting tools

### 🔧 Platform Features

#### User Management
- **Role-Based Access Control**: Separate interfaces for customers and companies
- **Secure Authentication**: Django's built-in authentication with custom user model
- **Profile Management**: Comprehensive user profile system
- **Email Verification**: Secure email-based account verification

#### Security & Performance
- **CSRF Protection**: Built-in CSRF protection for all forms
- **Input Validation**: Comprehensive form validation and sanitization
- **Error Handling**: Custom error pages (404, 403, 500)
- **Logging**: Comprehensive logging for debugging and monitoring

#### User Interface
- **Modern Design**: Clean, professional interface with modern UI/UX principles
- **Responsive Layout**: Mobile-first responsive design
- **Interactive Elements**: Hover effects, animations, and smooth transitions
- **Accessibility**: WCAG compliant design elements

## 🛠️ Technology Stack

### Backend
- **Framework**: Django 5.2.3
- **Language**: Python 3.8+
- **Database**: SQLite (Development) / PostgreSQL (Production)
- **Authentication**: Django's built-in authentication system
- **Forms**: Django Forms with custom validation

### Frontend
- **HTML5**: Semantic markup with modern standards
- **CSS3**: Custom styling with modern design principles
- **JavaScript**: Vanilla JS for interactive elements
- **Responsive Design**: Mobile-first approach

### Development Tools
- **Version Control**: Git
- **Package Management**: pip
- **Virtual Environment**: Python venv
- **Code Quality**: PEP 8 compliance

## 🏗️ Architecture

### Project Structure
```
netfix/
├── main/                    # Core application (homepage, navigation, errors)
│   ├── views.py            # Main views and error handlers
│   ├── urls.py             # Main URL routing
│   ├── middleware.py       # Custom middleware
│   └── templates/main/     # Core templates
├── users/                   # User management application
│   ├── models.py           # User, Customer, Company models
│   ├── views.py            # Authentication and profile views
│   ├── forms.py            # Registration and login forms
│   └── templates/users/    # User-related templates
├── services/                # Service management application
│   ├── models.py           # Service and ServiceRequest models
│   ├── views.py            # Service CRUD and request management
│   ├── forms.py            # Service creation and request forms
│   ├── services.py         # Business logic services
│   └── templates/services/ # Service-related templates
├── static/                  # Static assets
│   └── css/
│       ├── style.css       # Main stylesheet
│       └── fonts/          # Custom fonts
├── netfix/                  # Project configuration
│   ├── settings.py         # Django settings
│   ├── urls.py             # Main URL configuration
│   └── wsgi.py             # WSGI configuration
└── manage.py               # Django management script
```

### Application Architecture

#### Main App (`main/`)
- **Purpose**: Core application functionality
- **Responsibilities**: 
  - Homepage and navigation
  - Error handling (404, 403, 500)
  - Custom middleware
  - Base templates

#### Users App (`users/`)
- **Purpose**: User management and authentication
- **Responsibilities**:
  - User registration (Customer/Company)
  - Authentication and login
  - Profile management
  - Role-based access control

#### Services App (`services/`)
- **Purpose**: Service and request management
- **Responsibilities**:
  - Service creation and management
  - Service request handling
  - Request status management
  - Service analytics

## 🚀 Installation

### Prerequisites

- **Python**: 3.8 or higher
- **pip**: Python package installer
- **Git**: Version control system
- **Virtual Environment**: Python venv (recommended)

### Step-by-Step Installation

1. **Clone the Repository**
   ```bash
   git clone https://platform.zone01.gr/git/icelilog/netfix.git
   cd netfix
   ```

2. **Create Virtual Environment**
   ```bash
   # On Linux/macOS
   python3 -m venv venv
   source venv/bin/activate
   
   # On Windows
   python -m venv venv
   venv\Scripts\activate
   ```

3. **Install Dependencies**
   ```bash
   pip install django==5.2.3
   ```

4. **Configure Database**
   ```bash
   python manage.py makemigrations
   python manage.py migrate
   ```

5. **Run Development Server**
   ```bash
   python manage.py runserver
   ```

6. **Access Application**
   Open your browser and navigate to `http://127.0.0.1:8000`

## 📖 Usage

### For Customers

1. **Registration**
   - Visit the homepage and click "Register as Customer"
   - Fill in your details (name, email, birth date)
   - Complete registration and log in

2. **Browsing Services**
   - Navigate to "Services" to view all available services
   - Use category filters to find specific services
   - Click on services to view detailed information

3. **Requesting Services**
   - Select a service and click "Request Service"
   - Fill in service details (date, time, duration, address)
   - Submit request and wait for company response

4. **Managing Requests**
   - View your profile to see all service requests
   - Track request status and company responses
   - Access request history and details

### For Companies

1. **Registration**
   - Visit the homepage and click "Register as Company"
   - Fill in company details (name, email, field of work)
   - Complete registration and log in

2. **Creating Services**
   - Navigate to "Create Service" from your profile
   - Fill in service details (name, description, hourly rate, category)
   - Submit service for customer discovery

3. **Managing Requests**
   - View incoming service requests in your profile
   - Accept or reject requests with appropriate feedback
   - Track request status and customer information

4. **Company Profile**
   - Customize your public company profile
   - Showcase services and company information
   - Monitor business analytics and performance

## 🔌 API Endpoints

### Authentication Endpoints
- `POST /users/customer/` - Customer registration
- `POST /users/company/` - Company registration
- `POST /users/login/` - User login
- `GET /users/logout/` - User logout

### Service Endpoints
- `GET /services/` - List all services
- `GET /services/<id>/` - Get service details
- `POST /services/create/` - Create new service (companies only)
- `GET /services/most-requested/` - Get most requested services
- `GET /services/<field>/` - Get services by category

### Request Endpoints
- `POST /services/<id>/request/` - Request a service (customers only)
- `POST /services/request/<id>/<action>/` - Manage request (companies only)

### Profile Endpoints
- `GET /customer/<username>/` - Customer profile
- `GET /company/<username>/` - Company profile
- `GET /company-public/<username>/` - Public company page

### Service Categories
- Air Conditioner
- All in One
- Carpentry
- Electricity
- Gardening
- Home Machines
- House Keeping
- Interior Design
- Locks
- Painting
- Plumbing
- Water Heaters

## 🔒 Security Features

### Authentication & Authorization
- **Custom User Model**: Extended Django's AbstractUser for role-based access
- **Email-based Login**: Secure email/password authentication
- **Role-based Access Control**: Separate permissions for customers and companies
- **Session Management**: Secure session handling with CSRF protection

### Data Protection
- **Input Validation**: Comprehensive form validation and sanitization
- **SQL Injection Prevention**: Django ORM protection
- **XSS Protection**: Built-in Django security features
- **CSRF Protection**: Enabled for all forms and requests

### Error Handling
- **Custom Error Pages**: Professional 404, 403, and 500 error pages
- **Logging**: Comprehensive logging for security monitoring
- **Exception Handling**: Graceful error handling throughout the application


### Contribution Guidelines

- **Code Quality**: Ensure all code follows PEP 8 and Django conventions
- **Testing**: Add tests for new features
- **Documentation**: Update documentation for new features
- **Commit Messages**: Use clear, descriptive commit messages

### Version 1.0.0 (Current)
- Initial release
- Complete user management system
- Service creation and management
- Request handling system
- Modern responsive design
- Comprehensive security features

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👩‍💻 Authors

[Georgia Marouli](https://discordapp.com/users/1277216244910522371) - [Sofia Busho](https://discordapp.com/users/1276592724979613697)

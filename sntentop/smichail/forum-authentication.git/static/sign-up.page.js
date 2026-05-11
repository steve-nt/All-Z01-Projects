const form = document.getElementById('signup-form');
const successPopup = document.getElementById('success-popup');
const errorPopup = document.getElementById('error-popup');
const successMessage = document.getElementById('success-message');
const errorList = document.getElementById('error-list');

function showPopup(id) {
	document.getElementById(id).classList.add('show');
}

function closePopup(id) {
	document.getElementById(id).classList.remove('show');

	if (id === 'success-popup') {
		window.location.href = '/login';
	}
}

form.addEventListener('submit', async (e) => {
	e.preventDefault();

	// Clear messages
	errorList.innerHTML = '';

	const formData = {
		mail: document.getElementById('email-input').value,
		username: document.getElementById('username-input').value,
		password: document.getElementById('password-input').value,
		repeat_password: document.getElementById('repeat-password-input').value,
		role: 'user'
	};

	try {
		const response = await fetch('/signup', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(formData)
		});

		const data = await response.json();
		console.log('Response data:', data);

		if (data.success) {
			successMessage.textContent = data.message || 'Registration successful!';
			showPopup('success-popup');
			form.reset();
		} else {
			if (data.errors && data.errors.length > 0) {
				data.errors.forEach(error => {
					const li = document.createElement('li');
					li.textContent = error;
					errorList.appendChild(li);
				});
			} else {
				const li = document.createElement('li');
				li.textContent = data.message || 'An unknown error occurred.';
				errorList.appendChild(li);
			}
			showPopup('error-popup');
		}
	} catch (error) {
		const li = document.createElement('li');
		console.log(error);
		li.textContent = 'Network error. Please try again.';
		errorList.appendChild(li);
		showPopup('error-popup');
	}
});

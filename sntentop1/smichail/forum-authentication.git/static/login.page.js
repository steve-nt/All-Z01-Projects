
const form = document.getElementById('login-form');
const notVerifiedPopup = document.getElementById('not-verified-popup');
const notVerifiedMessage = document.getElementById('not-verified-message');
const resendEmailBtn = document.getElementById('resend-email-btn');

function showPopup(id) {
	document.getElementById(id).classList.add('show');
}
function closePopup(id) {
	document.getElementById(id).classList.remove('show');
}
let CurrerntUsername = '';

form.addEventListener('submit', async (e) => {
	e.preventDefault();

	const usernameInput = document.getElementById('username-input').value.trim();
	const passwordInput = document.getElementById('password-input').value;

	CurrerntUsername = usernameInput; // <-- αποθηκεύουμε το username

	const formData = { username: usernameInput, password: passwordInput };

	try {
		const response = await fetch('/login', {
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(formData)
		});


		const data = await response.json();
		console.log('Response data:', data);

		if (data.verified === true && data.success === true) {
			// Login successful, redirect to dashboard or home page
			window.location.href = '/posts';
		} else {
			// Login failed
			if (data.verified === false && data.resend === true) {
				// Show not verified popup
				notVerifiedMessage.textContent = data.message || 'Your email is not verified. Please verify your account to continue.';
				showPopup('not-verified-popup');
				form.reset();
			} else {
				alert(data.message || 'Login failed. Please try again.');
			}
		}
	} catch (error) {
		const li = document.createElement('li');
		li.textContent = 'Network error. Please try again.';
		errorList.appendChild(li);
		showPopup('not-verified-popup');
	}
});
document.addEventListener('DOMContentLoaded', () => {
	const resendEmailBtn = document.getElementById('resend-email-btn');

	resendEmailBtn.addEventListener('click', async () => {

		if (CurrerntUsername === '') {
			alert('Username is empty.');
			return;
		}

		try {
			const resendResponse = await fetch('/resend-verification', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ username: CurrerntUsername })
			});

			const resendData = await resendResponse.json();
			alert(resendData.message || 'Verification email resent!');
			document.getElementById('not-verified-popup').classList.remove('show');
		} catch (error) {
			console.error(error);
			alert('Network error while resending email.');
		}
	});
});


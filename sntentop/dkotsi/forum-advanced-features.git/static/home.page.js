const checkboxes = document.getElementsByName('categories');
const create_checkboxes = document.getElementsByName('create-categories');
const post_category = document.getElementsByName('post-category');
const button = document.getElementById('create-post-button');
const div = document.getElementById('create-post');
const likebtns = document.getElementsByName('like-btn');
const dislikebtns = document.getElementsByName('dislike-btn');
const deletebtns = document.getElementsByName('delete-btn-container');
const editbtns = document.getElementsByName('edit-btn-container');
const editforms = document.getElementsByName('edit-form')

function confirmLogout() {
	if (confirm("Are you sure you want to logout?")) {
		document.getElementById('logoutForm').submit();
	}
}
function showPopup(id) {
	document.getElementById(id).classList.add('show');
}

function closePopup(id) {
	document.getElementById(id).classList.remove('show');
}




// Load saved state from localStorage
checkboxes.forEach(checkbox => {
	const saved = localStorage.getItem(checkbox.id);
	if (saved !== null) {
		checkbox.checked = JSON.parse(saved);
	}

	// Save state to localStorage on change
	checkbox.addEventListener('change', () => {
		localStorage.setItem(checkbox.id, checkbox.checked);
	});
});


post_category.forEach(cat => {
	const checkboxid = cat.outerText + '-' + cat.id;
	localStorage.setItem(checkboxid, 'true');
});

create_checkboxes.forEach(checkbox => {
	console.log(checkbox.getAttribute("class"));
	if (checkbox.getAttribute("class") !== "create") {
		const saved = localStorage.getItem(checkbox.id);
		if (saved !== null) {
			checkbox.checked = JSON.parse(saved);
		}

		// Save state to localStorage on change
		checkbox.addEventListener('change', () => {
			localStorage.setItem(checkbox.id, checkbox.checked);
		});
	}
});


if (div) {
	div.style.display = 'none'; // Hides the div
	button.addEventListener('click', () => {
		if (div.style.display === 'none') {
			div.style.display = 'block'; // Resets to default
		} else {
			div.style.display = 'none'; // Hides the div
		}
	});
}


function reactonpostorcomment(btn, url) {

	btn.addEventListener('click', async (e) => {
		e.preventDefault();


		try {
			const response = await fetch(url + btn.id, {
				method: 'POST',
			});

			const data = await response.json();

			if (data.success) {
				console.log('success');
				window.location.reload();
			} else {
				showPopup('error-popup');
				console.log('error');
			}
		} catch (error) {
			showPopup('error-popup');
			console.log(error);
		}

	});
}
//add event listeners to react on click
likebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/like-post/');
});


dislikebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/dislike-post/');
});

deletebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/remove-post/');
});

editbtns.forEach((btn) => {
	console.log('hello');
	const article = document.getElementById('article-' + btn.id);
	const editform = document.getElementById('edit-form-' + btn.id);
	const closeditbtns = document.getElementById('close-edit-' + btn.id);
	editform.style.display = 'none';

	btn.addEventListener('click', function(e) {
		e.preventDefault();

		if (article.style.display === 'none') {
			article.style.display = 'flex'; // Resets to default
		} else {
			article.style.display = 'none'; // Hides the div
		}

		if (editform.style.display === 'none') {
			editform.style.display = 'flex'; // Resets to default
		} else {
			editform.style.display = 'none'; // Hides the div
		}
	});
	closeditbtns.addEventListener('click', () => {
		console.log('hello');
		editform.style.display = 'none';
		article.style.display = 'flex';
	});
});


editforms.forEach((form) => {

	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const formData = new FormData(form);
		console.log(form.action);

		console.log(formData);
		try {
			const response = await fetch(form.action, {
				method: 'POST',
				body: formData
			});

			const data = await response.json();

			if (data.success) {
				console.log('success');
				window.location.reload();
			} else {
				showPopup('error-popup');
				console.log('error');
			}
		} catch (error) {
			showPopup('error-popup');
			console.log(error);
		}
	});
});

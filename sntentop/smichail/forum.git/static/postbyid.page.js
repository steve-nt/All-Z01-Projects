const likebtns = document.getElementsByName('like-btn');
const dislikebtns = document.getElementsByName('dislike-btn');
const create_checkboxes = document.getElementsByName('create-categories');
const post_category = document.getElementsByName('post-category');
const commentdislikebtns = document.getElementsByName('comment-dislike-btn');
const commentlikebtns = document.getElementsByName('comment-like-btn');
const editpostbtns = document.getElementsByName('edit-post-btn-container');
const deletepostbtns = document.getElementsByName('delete-post-btn-container');
const editcommentbtns = document.getElementsByName('edit-comment-btn-container');
const deletecommentbtns = document.getElementsByName('delete-comment-btn-container');
const editpost = document.getElementById('edit-post');
const editcomment = document.getElementsByName('edit-comment-form');

post_category.forEach(cat => {
	const checkboxid = cat.outerText + '-' + cat.id;
	localStorage.setItem(checkboxid, 'true');
});

create_checkboxes.forEach(checkbox => {
	const saved = localStorage.getItem(checkbox.id);
	if (saved !== null) {
		checkbox.checked = JSON.parse(saved);
	}

	// Save state to localStorage on change
	checkbox.addEventListener('change', () => {
		localStorage.setItem(checkbox.id, checkbox.checked);
	});
});
function showPopup(id) {
	document.getElementById(id).classList.add('show');
}

function closePopup(id) {
	document.getElementById(id).classList.remove('show');
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

commentlikebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/like-comment/');
});

deletepostbtns.forEach((btn) => {
	reactonpostorcomment(btn, '/remove-post/');
});
dislikebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/dislike-post/');
});

commentdislikebtns.forEach((btn) => {
	reactonpostorcomment(btn, '/dislike-comment/');
});

deletecommentbtns.forEach((btn) => {
	reactonpostorcomment(btn, '/remove-comment/');
});

editcommentbtns.forEach((btn) => {
	console.log('hello');
	const edit_comment = document.getElementById('edit-comment-' + btn.id);
	const comment = document.getElementById('comment-' + btn.id);
	const closcommentedit = document.getElementById('close-edit-comment-button-container-' + btn.id);
	edit_comment.style.display = 'none';

	btn.addEventListener('click', function(e) {
		e.preventDefault();

		if (comment.style.display === 'none') {
			comment.style.display = 'flex'; // Resets to default
		} else {
			comment.style.display = 'none'; // Hides the div
		}

		if (edit_comment.style.display === 'none') {
			edit_comment.style.display = 'flex'; // Resets to default
		} else {
			edit_comment.style.display = 'none'; // Hides the div
		}
	});
	closcommentedit.addEventListener('click', () => {
		console.log('hello');
		edit_comment.style.display = 'none';
		comment.style.display = 'flex';
	});
});

editpostbtns.forEach((btn) => {
	console.log('hello');
	const edit_post = document.getElementById('edit-post');
	const post = document.getElementById('post');
	const clospostedit = document.getElementById('close-post-edit');
	edit_post.style.display = 'none';

	btn.addEventListener('click', function(e) {
		e.preventDefault();

		if (post.style.display === 'none') {
			post.style.display = 'flex'; // Resets to default
		} else {
			post.style.display = 'none'; // Hides the div
		}

		if (edit_post.style.display === 'none') {
			edit_post.style.display = 'flex'; // Resets to default
		} else {
			edit_post.style.display = 'none'; // Hides the div
		}
	});
	clospostedit.addEventListener('click', () => {
		console.log('hello');
		edit_post.style.display = 'none';
		post.style.display = 'flex';
	});
});

if (editpost) {
	editpost.addEventListener('submit', async (e) => {
		e.preventDefault();

		const formdata = new FormData(editpost);
		console.log(editpost.action);

		console.log(formdata);
		try {
			const response = await fetch(editpost.action, {
				method: 'post',
				body: formdata
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

editcomment.forEach((form) => {
	form.addEventListener('submit', async (e) => {
		e.preventDefault();
		const formdata = new FormData(form);
		console.log(form.action);

		try {
			const response = await fetch(form.action, {
				method: 'post',
				body: formdata
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

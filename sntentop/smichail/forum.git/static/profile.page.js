const notificationsbutton = document.getElementById('notification-button');

const notificationcontainer = document.getElementById('notifications');
const unseennotifications = document.getElementsByName('notification');
const created_section = document.getElementById('created-section');
const liked_section = document.getElementById('liked-section');
const activity_section = document.getElementById('activity-section');
const likedposts = document.getElementById('liked-posts');
const createdposts = document.getElementById('created-posts');
const activityelements = document.getElementById('activity-elements');
const sections = [created_section, liked_section, activity_section];
const elements = [createdposts, likedposts, activityelements];

sections.forEach((sec, index) => {
	sec.style.opacity = '0.4';
	elements[index].style.display = 'none';
});
const focused_section = localStorage.getItem('focus');
console.log(focused_section);
console.log('hello', sections[JSON.parse(focused_section)]);
if (focused_section !== null) {
	sections[JSON.parse(focused_section)].style.opacity = 1;
	elements[JSON.parse(focused_section)].style.display = 'flex';
} else {
	created_section.style.opacity = 1;
	createdposts.style.display = 'flex';
}

// liked_section.style.opacity = '0.4';
// activity_section.style.opacity = '0.4';
// created_section.style.opacity = '1';
notificationsbutton.addEventListener("click", () => {
	if (notificationcontainer.style.display === 'flex') {
		notificationcontainer.style.display = 'none'
	} else {
		notificationcontainer.style.display = 'flex'
	}
});

unseennotifications.forEach((not) => {
	const href = not.href
	not.addEventListener('click', async (e) => {
		e.preventDefault();


		try {
			const response = await fetch('/see-notification/' + not.id, {
				method: 'POST',
			});

			const data = await response.json();

			if (data.success) {
				console.log(data);
				window.location.href = href
			} else {
				console.log('error');
			}
		} catch (error) {
			console.log(error);
		}
	});
});




sections.forEach((section, index) => {
	section.addEventListener('click', () => {

		if (section.style.opacity !== '1') {
			localStorage.setItem('focus', index);
			sections.forEach((sec) => {
				sec.style.opacity = '0.4';
			});
			elements.forEach((element) => {
				element.style.display = 'none';
			});
			section.style.opacity = '1';
			switch (section.getAttribute('name')) {
				case "liked":
					likedposts.style.display = 'flex';
					break;
				case "created":
					createdposts.style.display = 'flex';
					break;
				case "activity":
					activityelements.style.display = 'flex';
					break;
			};
		}
	});
});


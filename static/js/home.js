function elem(name) {
	let element = document.getElementById(name);
	if (element == null) {
		element = document.getElementsByClassName(name)[0];
	}

	return {
		addclass: function(clsname) {
			element.classList.add(clsname);
			return this;
		},
		removeclass: function(clsname) {
			element.classList.remove(clsname);
			return this;
		},
		mouseover: function(fn) {
			element.addEventListener('mouseover', fn);
			return this;
		}
	};
}

function set_navbar_color(pos) {
	if (pos >= 200)
		$('.navbar').addClass('bg-dark');
	else
		$('.navbar').removeClass('bg-dark');
}

(() => {
	$('.navbar').addClass('fixed-top')
	let overlay = $('.overlay');
	let default_bg = overlay.css('background-color');

	$(window).scroll(() => {
		let pos = $(window).scrollTop();

		set_navbar_color(pos);

		if (pos >= 550)
			overlay.css('background-color', 'black');
		else
			overlay.css('background-color', default_bg);
	})
	elem('arrow').mouseover(() => {});
})();
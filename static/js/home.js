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

(() => {
	let nav = elem('nav').addclass('fixed-top');

	this.onscroll = () => {
		if (pageYOffset >= 200)
			nav.addclass('bg-dark');
		else
			nav.removeclass('bg-dark');
	}
	// elem('arrow').mouseover(() => {});
})();
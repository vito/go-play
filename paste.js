$(function(){
	$(".paste-input").keydown(function(e) {
		if (e.keyCode != 9 || e.ctrlKey || e.altKey)
			return;

		if (this.setSelectionRange) {
			var start = this.selectionStart,
				end = this.selectionEnd,
				top = this.scrollTop;

			this.value = this.value.slice(0, start) + '\t' +
			this.value.slice(end);
			this.setSelectionRange(start + 1, start + 1);
			this.scrollTop = top;

			e.preventDefault();
		} else if (document.selection.createRange) {
			this.selection = document.selection.createRange();
			this.selection.text = '\t';
			e.returnValue = false;
		}	
	});
});

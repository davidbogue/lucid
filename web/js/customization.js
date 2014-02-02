$.fn.extend({
    insertAtCaret: function(myValue) {
        var txtArea = this[0];
        if (txtArea.selectionStart || txtArea.selectionStart == '0') {
            var startPos = txtArea.selectionStart;
            var endPos = txtArea.selectionEnd;
            var scrollTop = txtArea.scrollTop;
            txtArea.value = txtArea.value.substring(0, startPos)+myValue+txtArea.value.substring(endPos,txtArea.value.length);
            txtArea.focus();
            txtArea.selectionStart = startPos + myValue.length;
            txtArea.selectionEnd = startPos + myValue.length;
            txtArea.scrollTop = scrollTop;
        } else {
            txtArea.value += myValue;
            txtArea.focus();
        }
    }
})
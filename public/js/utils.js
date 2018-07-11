function minimalDelay(min) {
    var lastTimestamp = Date.now(); // in ms
    return function (fn) {
        var passed = Date.now() - lastTimestamp;
        if (passed < min) {
            var left = min - passed;
            setTimeout(function () {
                fn();
            }.bind(this), left);
        } else {
            fn();
        }
    }
}

function fromCache(vue) {

    var cached = {};

    try {
        cached = JSON.parse(window.localStorage.getItem('form'))
    } catch (e) {
    }

    for (var key in cached) {
        if (vue.form.hasOwnProperty(key)) {
            vue.form[key] = cached[key];
        }
    }
}

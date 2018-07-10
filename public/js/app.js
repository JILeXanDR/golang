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

var App = Vue.component('app', {
    template: '#template_app',
    data: function () {
        return {}
    }
});

Vue.component('order-form', {
    template: '#template_order-form',
    data: function () {
        return {
            item: '',
            form: {
                list: [],
                phone: '',
                delivery_address: '',
                first_name: '',
                comment: '',
            },
            loader: false,
            search: '',
            dialog: {
                visible: false,
                text: '',
            },
            items: [],
            isLoading: false,
        }
    },
    created() {
        this.sync();
        this.search = this.form.delivery_address;
    },
    methods: {
        add() {
            this.form.list.push(this.item);
            this.item = '';
        },
        remove(item) {
            var index = this.form.list.indexOf(item);
            this.form.list.splice(index, 1);
        },
        confirm() {

            this.loader = true;

            var wait = minimalDelay(500);

            this.$http.post('/order', this.form).then(function (res) {
                wait(function () {
                    this.loader = false;
                    this.dialog.text = 'Заказ оформлен. Ожидайте доставку курьером';
                    this.dialog.visible = true;
                }.bind(this));
            }).catch(function (err) {
                wait(function () {
                    this.loader = false;
                    this.dialog.text = err.body;
                    this.dialog.visible = true;
                }.bind(this));
            });
        },
        sync() {

            var cached = {};

            try {
                cached = JSON.parse(window.localStorage.getItem('form'))
            } catch (e) {
            }

            for (var key in cached) {
                if (this.form.hasOwnProperty(key)) {
                    this.form[key] = cached[key];
                }
            }
        },
        isFormValid() {
            return this.form.list.length && this.form.phone && this.form.first_name && this.form.delivery_address;
        }
    },
    watch: {
        form: {
            handler: function () {
                window.localStorage.setItem('form', JSON.stringify(this.form));
            },
            deep: true
        },
        search(val) {

            this.isLoading = true;

            this.$http.get('/find-address', {params: {q: val}})
                .then(res => {
                    console.log(res.data);
                    this.items = res.data;
                })
                .catch(err => {
                    console.log(err)
                })
                .finally(() => (this.isLoading = false))
        }
    }
});

new Vue({
    el: '#app',
    render: function (createElement) {
        return createElement(App)
    }
});

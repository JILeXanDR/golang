Vue.component('home', {
    template: '#template_home',
    data: function () {
        return {
            item: '',
            form: {
                list: [],
                phone: '',
                delivery_address: '',
                name: '',
                comment: '',
            },
            loader: false,
            search: '',
            dialog: {
                visible: false,
                text: '',
            },
            dialogConfirmPhone: {
                visible: false,
                text: '',
            },
            items: [],
            isLoading: false,
            phoneConfirmed: false,
            smsCode: '',
            md5: '',
            snackbar: {
                visible: false,
                text: '',
                timeout: 3000,
            },
        }
    },
    created() {
        fromCache(this);
        if (this.form.delivery_address) {
            this.search = this.form.delivery_address.name;
        }
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

            var self = this;

            var showDialog = function (text) {
                this.loader = false;
                this.dialog.text = text;
                this.dialog.visible = true;
            }.bind(this);

            if (!this.phoneConfirmed) {
                this.dialogConfirmPhone.visible = true;
                return;
            }

            var waiter = minimalDelay(500);

            this.loader = true;

            this.$http.post('/api/orders', this.form).then(function (res) {
                waiter(function () {
                    showDialog('Заказ оформлен. Ожидайте СМС с подтверждением заказа');
                    var unwatch = self.$watch('dialog.visible', function (newValue, oldValue) {
                        // dialog was closed
                        if (!newValue) {
                            unwatch();
                            self.form.list = []; // TODO не работает
                            self.$router.push({name: 'orders'})
                        }
                    });
                });
            }).catch(function (err) {
                waiter(function () {
                    this.error(err.body.message);
                    showDialog(err.body.message);
                });
            });
        },
        isFormValid() {
            return this.form.list.length && this.form.phone && this.form.name && this.form.delivery_address;
        },
        checkSmsCode(code) {
            this.$http.post('/api/check-code', {code})
                .then(res => {
                    this.phoneConfirmed = true;
                    this.dialogConfirmPhone.visible = false;
                    this.confirm();
                })
                .catch(err => {
                    this.error(err.body.message);
                });
        },
        error(text) {
            this.snackbar.visible = true;
            this.snackbar.text = text;
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
                    this.items = res.data.map(function (item) {
                        return {
                            value: item.id,
                            name: item.name,
                        };
                    });
                })
                .catch(err => {
                    this.error(err.body.message);
                })
                .finally(function () {
                    this.isLoading = false;
                })
        },
        'dialogConfirmPhone.visible': function (val) {
            if (val) {
                this.$http.post('/api/confirm-phone', {phone: this.form.phone})
                    .then(res => {
                        this.md5 = res.body.md5;
                    })
                    .catch(err => {
                        this.error(err.body.message);
                    })
            } else {
                this.md5 = this.smsCode = '';
            }
        },
        smsCode: function (val) {
            // TODO use double hash
            if (this.md5 === md5(val)) {
                this.checkSmsCode(val);
            }
        }
    }
});

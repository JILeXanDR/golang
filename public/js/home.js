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
            // phoneConfirmed: false,
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
        this.$on('orderCreated', function () {
            this.form.list = []; // TODO не работает
            this.$router.push({name: 'orders'})
        });
    },
    methods: {
        add() {
            this.form.list.push(this.item);
            this.item = '';
        },
        remove(item) {
            this.form.list.splice(this.form.list.indexOf(item), 1);
        },
        confirmOrder() {

            let showDialog = (text) => {
                this.loader = false;
                this.dialog.text = text;
                this.dialog.visible = true;
            };

            // if (!this.phoneConfirmed) {
            //     this.dialogConfirmPhone.visible = true;
            //     return;
            // }

            let waiter = minimalDelay(500);

            this.loader = true;

            this.$http.post('/api/orders', this.form).then(res => waiter(() => {
                showDialog('Заказ оформлен. Ожидайте СМС с подтверждением заказа');
                let unwatch = this.$watch('dialog.visible', (newValue, oldValue) => {
                    // dialog was closed
                    if (!newValue) {
                        unwatch();
                        this.$emit('orderCreated');
                    }
                });
            })).catch((err) => waiter(() => {
                if (err.status === 401) {
                    this.dialogConfirmPhone.visible = true;
                } else {
                    this.error(err.body.message);
                    // showDialog(err.body.message);
                }
            }));
        },
        isFormValid() {
            return this.form.list.length && this.form.phone && this.form.name && this.form.delivery_address;
        },
        checkSmsCode(code) {
            this.$http.post('/api/check-code', {code})
                .then(res => {
                    // this.phoneConfirmed = true;
                    this.dialogConfirmPhone.visible = false;
                    this.confirmOrder();
                })
                .catch(err => this.error(err.body.message));
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
            this.$http.get('/find-address', {params: {q: val}})
                .then(res => {
                    this.items = res.data.map(function (item) {
                        return {
                            value: item.id,
                            name: item.name,
                        };
                    });
                })
                .catch(err => this.error(err.body.message));
        },
        'dialogConfirmPhone.visible': function (val) {
            if (val) {
                this.$http.post('/api/confirm-phone', {phone: this.form.phone})
                    .then(res => {
                        this.md5 = res.body.md5;
                        setTimeout(() => {
                            // TODO for fast dev
                            if (res.body.code) {
                                this.smsCode = String(res.body.code);
                            }
                        }, 1000);
                    })
                    .catch(err => this.error(err.body.message))
            } else {
                this.md5 = this.smsCode = '';
            }
        },
        smsCode: function (val) {
            if (this.md5 === md5(md5(val))) {
                this.checkSmsCode(val);
            }
        }
    },
});

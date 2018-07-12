Vue.component('order-form', {
    template: '#template_order-form',
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
            items: [],
            isLoading: false,
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

            var showDialog = function (text) {
                this.loader = false;
                this.dialog.text = text;
                this.dialog.visible = true;
            }.bind(this);

            var waiter = minimalDelay(500);

            this.loader = true;

            this.$http.post('/api/orders', this.form).then(function (res) {
                waiter(function () {
                    showDialog('Заказ оформлен. Ожидайте СМС с подтверждением заказа');
                });
            }).catch(function (err) {
                waiter(function () {
                    showDialog(err.body.message);
                });
            });
        },
        isFormValid() {
            return this.form.list.length && this.form.phone && this.form.name && this.form.delivery_address;
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

            // if (this.form.delivery_address && val === null) {
            //     return;
            // }

            // console.log(this.form.delivery_address);

            // if (typeof val === "object" && val !== null && val.value) {
            //     return;
            // }

            // console.log({val, type: typeof val});

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
                    console.log(err)
                })
                .finally(function () {
                    this.isLoading = false;
                })
        }
    }
});

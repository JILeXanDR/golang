Vue.component('orders', {
    template: '#tpl_orders',
    data: function () {
        return {
            headers: [
                {
                    text: 'ID',
                    align: 'left',
                    sortable: false,
                    value: 'id'
                },
                {
                    text: 'Покупки',
                    align: 'left',
                    sortable: false,
                    value: 'list'
                },
                {
                    text: 'Статус',
                    align: 'left',
                    sortable: false,
                    value: 'status'
                },
                {
                    text: 'Номер телефона',
                    align: 'left',
                    sortable: false,
                    value: 'phone'
                },
                {
                    text: 'Адрес доставки',
                    align: 'left',
                    sortable: false,
                    value: 'delivery_address'
                },
                {
                    text: 'Примечание',
                    align: 'left',
                    sortable: false,
                    value: 'comment'
                },
                {
                    text: 'Создан',
                    align: 'left',
                    sortable: false,
                    value: 'created_at'
                },
            ],
            orders: [],
        }
    },
    created() {
        this.$http.get('/api/orders')
            .then(res => {
                this.orders = res.body;
            })
            .catch(err => {

            })
    },
});

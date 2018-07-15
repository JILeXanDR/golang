Vue.component('orders-table', {
    template: '#tpl_orders-table',
    props: {
        data: {
            type: Array,
            required: true,
        },
    },
    data: function () {
        return {
            headers: [
                {
                    text: 'Покупки',
                    align: 'left',
                    sortable: false,
                    value: 'list'
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
                    text: 'Статус',
                    align: 'left',
                    sortable: false,
                    value: 'status'
                },
                {
                    text: 'Создан',
                    align: 'left',
                    sortable: false,
                    value: 'created_at'
                },
            ],
        }
    },
});


Vue.component('orders', {
    template: '#tpl_orders',
    data: function () {
        return {
            filterStatus: 'new',
            tabs: [
                {
                    id: 'new',
                    name: 'Новые',
                },
                {
                    id: 'done',
                    name: 'Выполненные',
                }
            ],
            orders: [],
        }
    },
    computed: {
        newOrders() {
            return this.orders.filter((order) => ['created', 'confirmed', 'processing'].includes(order.status));
        },
        doneOrders() {
            return this.orders.filter((order) => ['canceled', 'delivered'].includes(order.status));
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

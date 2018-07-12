const routes = [
    {
        path: '/',
        name: 'home',
        component: {template: '<order-form></order-form>'},
    },
    {
        path: '/orders',
        name: 'orders',
        component: {template: '<orders></orders>'},
    },
    {
        path: '/about-us',
        name: 'about_us',
        component: {template: '<about-us></about-us>'},
    },
];

const router = new VueRouter({
    mode: 'history',
    routes
});

Vue.filter('date', function (date) {
    var d = new Date(date);
    return `${d.getFullYear()}-${d.getMonth() + 1}-${d.getDate()}`;
});

new Vue({
    el: '#app',
    router,
    render: function (createElement) {
        return createElement(App)
    }
});

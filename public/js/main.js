const routes = [
    {
        path: '/',
        name: 'home',
        component: {template: '<home></home>'},
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

moment.locale('ru');

Vue.filter('date', function (date) {
    return moment(date).format("dddd, MMMM Do YYYY, H:mm:ss");
});

Vue.filter('timeLeft', function (date) {
    return moment(date).from(moment());
});

Vue.filter('phone', function (phone) {
    return phone.replace(/^(38)+/, '');
    return phone.replace(/[^0-9]/g, '').replace(/(\d{3})(\d{2})(\d{2})(\d{3})/, '($1) $2-$2-$3');
});

new Vue({
    el: '#app',
    router,
    render: function (createElement) {
        return createElement(App)
    }
});

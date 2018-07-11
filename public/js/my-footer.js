Vue.component('my-footer', {
    template: '#tpl_footer',
    data: function () {
        return {
            year: (new Date()).getFullYear(),
            siteName: 'SimpleShopping',
            icons: [
                'fab fa-instagram',
                'fab fa-facebook-square',
                'fab fa-vk',
            ],
        }
    }
});

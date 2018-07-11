Vue.component('my-footer', {
    template: '#tpl_footer',
    data: function () {
        return {
            year: (new Date()).getFullYear(),
            siteName: 'SimpleShopping',
            icons: [
                'delete',
                'delete',
                'delete',
                'delete',
                'delete',
            ],
        }
    }
});

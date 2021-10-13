(function ($) {
    'use strict';

    var Selector = {
        init: function () {
            $(document).ready(function () {
                const paths = window.location.pathname.split("/").filter(p => p)
                if (paths.length === 0) {
                    return
                }
                const version = paths[0].replace('docs', '').replace("-", '');
                if (version) {
                    $('.navbar-version-menu')[0].text = version;
                };
            });
        },
    };

    Selector.init();
}(jQuery));

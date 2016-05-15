
require.config({
        //By default load any module IDs from scripts/app
        baseUrl: (_baseUrl ||'') + '/assets/js/',
        //except, if the module ID starts with "lib"
        paths: {
            shop: 'touch/shop', //以shop开头的前缀,从那个路径找文件
            uc : 'touch/uc',
            'jquery':'lib/jquery.2x',
            'jquery.slides':'lib/jquery.slides',
        },
        // load backbone as a shim
        shim: { //依赖关系
            'backbone': {
                //The underscore script dependency should be loaded before loading backbone.js
                deps: ['underscore'],
                // use the global 'Backbone' as the module name.
                exports: 'Backbone'
            },
            'shop/main': {
                //The underscore script dependency should be loaded before loading backbone.js
                deps: ['jr/core'],
                // use the global 'Backbone' as the module name.
                exports: 'Main'
            },
            'uc/main':{
                deps :['jr/core']
            },
            'jr/scroller':{
                deps:['jr/core']
            },
            'jr/ui.min':{
                deps:['jr/core']
            },
            'jr/dialog':{
                deps:['jr/core']
            },
            'jquery.slides':{
                deps:['jquery']
            }
        }
    }
);
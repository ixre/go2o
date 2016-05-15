
require.config({
        //By default load any module IDs from scripts/app
        baseUrl: (_baseUrl ||'') + '/assets/js/',
        //except, if the module ID starts with "lib"
        paths: {
            shop: 'touch/shop', //以shop开头的前缀,从那个路径找文件
            'mui':'../mui/js',
            'jquery':'../easyui/jquery.min',
            'jquery.easyui.min':'../easyui/jquery.easyui.min',
            'jquery.easyui':'../easyui/locale/easyui-lang-zh_CN',
        },
        // load backbone as a shim
        shim: { //依赖关系
            'bm/main': {
                //The underscore script dependency should be loaded before loading backbone.js
                deps: ['jr/core'],
                // use the global 'Backbone' as the module name.
                exports: 'Main'
            },
            'jquery.easyui.min':{
                deps:['jquery']
            },
            'jquery.easyui':{
                deps:['jquery.easyui.min']
            },
            'mui/component':{
                deps:['jr/core']
            }
        }
    }
);
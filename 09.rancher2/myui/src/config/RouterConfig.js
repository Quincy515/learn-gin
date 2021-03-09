import VueRouter from "vue-router"

const RouterConfig = {
    routes: [
        { path: "/", name: "index", meta: { navname: "首页", index: 1 }, component: () => import(/* webpackChunkName: "index-chunk" */"@/components/index") },
        { path: "/login", name: "login", meta: { navname: "登录", index: 2 }, component: () => import(/* webpackChunkName: "login-chunk" */"@/components/login") },

    ]
};

export default new VueRouter(RouterConfig)
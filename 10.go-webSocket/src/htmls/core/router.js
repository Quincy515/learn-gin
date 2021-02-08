const routes = [
    { path: '/', component: httpVueLoader( 'components/table.vue' ) },
];
const router = new VueRouter({
    routes
});
export default router
import { Layout } from '@/routers/constant';

const appStoreRouter = {
    sort: 2,
    path: '/apps',
    component: Layout,
    redirect: '/apps/all',
    meta: {
        icon: 'p-appstore',
        title: '应用管理',
    },
    children: [
        {
            path: '/apps',
            name: 'App',
            redirect: '/apps/jar',
            meta: {},
            children: [
                {
                    path: 'jar',
                    name: 'Jar',
                    component: () => import('@/views/app-manage/jar/index.vue'),
                    props: true,
                    hidden: true,
                    meta: {
                        activeMenu: '/apps',
                        requiresAuth: false,
                    },
                },
            ],
        },
    ],
};

export default appStoreRouter;

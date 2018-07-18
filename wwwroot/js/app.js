'use strict';
const EventRegistered = 'REGISTERED';
const EventLoggedIn = 'LOGIN';
const EventLoggedOut = 'LOGOUT';
const EventBus = new Vue();
const appTitle = 'Template';

const navbar_component = {
    props: ['user'],
    template: `<div class="navbar navbar-expand-md bg-light fixed-top">
    <div class="container-fluid">
        <button type="button" class="navbar-toggler" data-toggle="collapse" data-target=".navbar-collapse">
            <i class="fas fa-bars"></i>
        </button>
        <router-link to="/" class="navbar-brand" title="` + appTitle + '">' + appTitle + `</router-link>
        <div class="navbar-collapse collapse">
            <router-link to="/" class="nav-link" title="Home">Home <i class="fas fa-home"></i></router-link>
            <router-link to="/account/register" v-if="!user" class="nav-link ml-auto" title="Register">Register</router-link>
            <router-link to="/account/signin" v-if="!user" class="nav-link" title="Login">Login</router-link>
            <div class="nav-item dropdown ml-auto" v-if="user">
                <a class="nav-link dropdown-toggle" href="#" id="navbarDropdownMenuLink" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                    Menu
                </a>
                <div class="dropdown-menu dropdown-menu-right" aria-labelledby="navbarDropdownMenuLink">
                    <p class="dropdown-item disabled" v-if="user">{{user.username}}</p>
                    <router-link to="/account/manage" class="dropdown-item" title="Manage">Settings</router-link>
                    <router-link to="/account/logout" class="dropdown-item text-danger" title="Logout">Logout</router-link>
                </div>
            </div>
        </div>
    </div>
</div>`
};
const register_component = {
    data() {
        return {
            username: {
                value: '',
                pattern: /^[A-Za-z0-9]+(?:[_-][A-Za-z0-9]+)*$/,
                error: 'Username can only contain \'a-Z, 0-9, -, _\''
            },
            password: {
                value: '',
                valid() {
                    return this.value.length > 5;
                },
                error: 'Password must be atleast 6 characters long'
            },
            password2: {
                value: '',
                error: 'passwords do not match'
            },
            email: {
                value: '',
                pattern: /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
                error: 'Email must be like \'myemail@provider.com\''
            },
            errors: {
                password: false,
                password2: false,
                username: false,
                email: false
            }
        };
    },
    template: `<div>
    <h2>Register</h2>
    <div class="row">
        <div class="col-md-6 col-lg-4">
            <h4>Create a new account.</h4>
            <div class="form-group">
                <label>Username</label>
                <input required v-model="username.value" class="form-control" />
                <span v-if="errors.username" class="text-danger">{{username.error}}</span>
            </div>
            <div class="form-group">
                <label>Email</label>
                <input required v-model="email.value" class="form-control" />
                <span v-if="errors.email" class="text-danger">{{email.error}}</span>
            </div>
            <div class="form-group">
                <label>Password</label>
                <input required v-model="password.value" type="password" class="form-control" />
                <span v-if="errors.password" class="text-danger">{{password.error}}</span>
            </div>
            <div class="form-group">
                <label>Confirm password</label>
                <input required v-model="password2.value" type="password" class="form-control" />
                <span v-if="errors.password2" class="text-danger">{{password2.error}}</span>
            </div>
            <button class="btn btn-default" @click="register_new" >Register</button>
        </div>
    </div>
</div>`,
    methods: {
        validate_username() {
            return !(this.errors.username = !(new RegExp(this.username.pattern).test(this.username.value)));
        },
        validate_email() {
            return !(this.errors.email = !(new RegExp(this.email.pattern).test(this.email.value)));
        },
        validate_password() {
            return !(this.errors.password = !this.password.valid());
        },
        validate_password2() {
            return !(this.errors.password2 = (this.password.value !== this.password2.value));
        },
        register_new() {
            const valid = this.validate_username() & this.validate_email() &
                this.validate_password() & this.validate_password2();
            if (valid) {
                const payload = {
                    username: this.username.value,
                    password: this.password.value,
                    email: this.email.value
                };
                this.$signInManager.Register(payload)
                    .then(() => {})
                    .catch(() => {});
            }
        }
    }
};
const login_component = {
    props: ['external_logins'],
    data() {
        return {
            errors: {
                request: false
            },
            username: '',
            password: '',
            remember_me: false
        };
    },
    template: `<div>
    <h2>Log in</h2>
    <div class="row">
        <div class="col-md-6 col-lg-4">
            <h4>Log in using your account</h4>
            <hr />
            <form>
                <div v-if="errors.request && errors.request !== ''" class="alert alert-danger" role="alert">{{errors.request}}</div>
                <div class="form-group">
                    <label>Username</label>
                    <input required v-model="username" class="form-control" />
                </div>
                <div class="form-group">
                    <label>Password</label>
                    <input required v-model="password" type="password" class="form-control" />
                </div>
                <label>Remember me
                    <input v-model="remember_me" type="checkbox" />
                </label>
                <div class="form-group">
                    <button @click="login" type="submit" class="btn btn-default">Log in</button>
                </div>
            </form>
            <p>
                <a asp-page="./ForgotPassword">Forgot your password?</a>
            </p>
            <p>
                <router-link to="/account/register">Register as a new user</router-link>
            </p>
        </div>
        <div class="col-md-6 col-md-offset-2">
        <div v-if="external_logins && external_logins.length > 0">
            <h4>Use another service to log in.</h4>
            <hr />
                <div>
                    <p>
                        There are no external authentication services configured. See <a href="https://go.microsoft.com/fwlink/?LinkID=532715">this article</a>
                        for details on setting up this ASP.NET application to support logging in via external services.
                    </p>
                </div>
                <div>
                    <p>
                        <button v-for="e in external_logins" type="submit" class="btn btn-default" name="provider" value="{{e.name}}" title="Log in using your {{e.display_name}} account">{{e.display_name}}</button>
                    </p>
                </div>
            </div>
        </div>
    </div>`,
    methods: {
        login() {
            const vm = this;
            this.$signInManager.SignIn(vm.username, vm.password, vm.remember_me)
                .catch((err) => {
                    vm.errors.request = err.response.data;
                });
        }
    }
};
const logout_component = {
    data() {
        return {
            status: 'logging out...'
        };
    },
    template: '<h2>{{status}}</h2>',
    created() {
        const vm = this;
        this.$signInManager.LogOut().then(() => {
            vm.status = 'Logged out';
        });
    }
};

const home_component = {
    props: ['authenticated'],
    template: `<div>
    Hello World!
</div>`
};

Vue.use(VueRouter);

const routes = [{
    path: '/',
    component: home_component
}, {
    path: '/account/register',
    component: register_component
}, {
    path: '/account/signin',
    component: login_component
}, {
    path: '/account/logout',
    component: logout_component
}];
const router = new VueRouter({
    routes
});

class SignInManager {
    constructor(http) {
        this.http = http;
        this.user = this.GetUser();
    }
    SignIn(username, password, remember) {
        const sim = this;
        return new Promise((res, rej) => {
            if (username.length > 1 && password.length > 5) {
                const payload = {
                    username,
                    password,
                    grant_type: 'password',
                };
                var queryString = Object.keys(payload).map(function (key) {
                    return encodeURIComponent(key) + '=' + encodeURIComponent(payload[key]);
                }).join('&');
                sim.http({
                    method: 'post',
                    url: '/api/token',
                    data: queryString,
                    headers: {
                        'Content-type': 'application/x-www-form-urlencoded'
                    }
                }).then((data) => {
                    const token = data.data.access_token;
                    if (remember) {
                        localStorage.setItem('token', token);
                    } else {
                        sessionStorage.setItem('token', token);
                    }
                    sim.user = sim.parseJwt(token);
                    EventBus.$emit(EventLoggedIn);
                    res();
                }).catch(rej);
            } else {
                rej('Please enter a valid username/password');
            }
        });
    }
    LogOut() {
        return new Promise((res) => {
            localStorage.removeItem('token');
            sessionStorage.removeItem('token');
            this.user = null;
            EventBus.$emit(EventLoggedOut);
            res();
        });
    }
    Register(user) {
        const sim = this;
        return new Promise((res, rej) => {
            sim.http.post('/account/register', user)
                .then(() => {
                    EventBus.$emit(EventRegistered);
                    res();
                })
                .catch(rej);
        });
    }
    GetUser() {
        const token = this.GetToken();
        if (!token) return null;

        return this.parseJwt(token);
    }
    GetLogOutTime() {
        if (!this.user) return null;
        return new Date(this.user.exp * 1000);
    }
    LoggedIn() {
        return (localStorage.getItem('token') !== null) ||
            (sessionStorage.getItem('token') !== null);
    }
    GetToken() {
        return localStorage.getItem('token') || sessionStorage.getItem('token');
    }
    parseJwt(token) {
        var base64Url = token.split('.')[1];
        var base64 = base64Url.replace('-', '+').replace('_', '/');
        return JSON.parse(atob(base64));
    }

}

Vue.prototype.$http = axios;
Vue.prototype.$signInManager = new SignInManager(axios);
new Vue({
    data() {
        return {
            user: undefined
        };
    },
    created() {
        const vm = this;
        vm.user = vm.$signInManager.user;
        EventBus.$on(EventLoggedIn, function () {
            vm.user = vm.$signInManager.user;
            vm.$router.push('/');
        });
        EventBus.$on(EventLoggedOut, function () {
            vm.user = vm.$signInManager.user;
            vm.$router.push('/account/signin');
        });
        EventBus.$on(EventRegistered, function () {
            vm.$router.push('/account/signin');
        });
    },
    router,
    el: '#app',
    template: `<div class="container">
    <navbar v-bind:user="user"></navbar>
    <router-view></router-view>
</div>`,
    components: {
        'navbar': navbar_component
    }
});
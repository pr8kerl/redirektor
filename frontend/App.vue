<template>
  <div id="app">

   <header>
     <nav class="grey darken-2" role="navigation">
       <div class="nav-wrapper container">
         <a href="#" class="brand-logo left" >
           <i class="material-icons">autorenew</i>
         </a>
         <ul id="willkommen" class="right">
           <li>welcome luser</li>
         </ul>
       </div>
     </nav>
   </header>

    <v-container>
      <div class="section no-pad-bot center" id="index-banner">
        <img src="/assets/img/logo.png">
        <h1 class="header center orange-text">{{ msg }}</h1>

        <!--[if lt IE 8]>
        <div class="row center pink-text"><h2>You are using an <strong>outdated</strong> browser. Please <a href="http://browsehappy.com/">upgrade your browser</a> to improve your experience.</h2></div>
        <![endif]-->

          <div v-if="loading">
            <v-progress-circular active red red-flash></v-progress-circular>
          </div>

          <div v-if="error" class="error row center red-text">
            {{ error }}
          </div>

          <div v-if="redirekts" class="content">
            <input v-model="filter" class="form-control" placeholder="search by incoming">
            <table>
              <thead>
                <tr>
                    <th data-field="prefix">prefix</th>
                    <th data-field="incoming">incoming</th>
                    <th data-field="outgoing">outgoing</th>
                </tr>
              </thead>
              <tbody>
              <tr v-for="redirekt in filterRedirekts(filter)">
              <td>{{ redirekt.prefix }}</td>
              <td>{{ redirekt.incoming }}</td>
              <td>{{ redirekt.outgoing }}</td>
              </tr>
              </tbody>
            </table>
          </div>

      </div>
    </v-container>
  </div>
</template>

<script>
export default {
  name: 'app',
  data () {
    return {
      msg: 'Welcome to Your Vue.js App',
      redirekts: null,
      loading: true,
      filter: '',
      error: null
    }
  },
  methods: {
    redirektsOk: function(response) {
          this.redirekts = response.body.response;
          this.loading = false;
          console.log(this.redirekts);
    },
    redirektsError: function(response) {
          this.loading = false;
          this.error = 'there was an error';
          console.log(JSON.stringify(response.statusText));
    },
    filterRedirekts: function(value) {
            return this.redirekts.filter(function(item) {
                return item.prefix.indexOf(value) > -1 ||
                       item.incoming.indexOf(value) > -1;
            });
    }
  },
  created: function () {
        // GET redirekts request
        this.$http.get('/api/redirekts').then(this.redirektsOk,this.redirektsError);
  }
}
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

h1, h2 {
  font-weight: normal;
}

ul {
  list-style-type: none;
  padding: 0;
}

li {
  display: inline-block;
  margin: 0 10px;
}

a {
  color: #42b983;
}
</style>

import Editor from './Editor.vue'

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

    <div>
      <div class="section no-pad-bot center" id="index-banner">
        <h2 class="header center orange-text">{{ msg }}</h2>

        <!--[if lt IE 8]>
        <div class="row center pink-text"><h2>You are using an <strong>outdated</strong> browser. Please <a href="http://browsehappy.com/">upgrade your browser</a> to improve your experience.</h2></div>
        <![endif]-->

          <div v-if="loading">
            <v-progress-circular active red red-flash></v-progress-circular>
          </div>

          <div v-if="error" class="error row center red-text">
            {{ error }}
          </div>


          <div v-if="redirektdata" class="content">
            <div class="row">
              <div class="col s1">
                <v-icon prefix>search</v-icon>
              </div>
              <div class="col s10">
                <input v-model="filter" class="form-control" placeholder="filter incoming">
              </div>
            </div>
            <div class="row">
            <table id="tableredirekts" class="striped">
              <thead>
                <tr>
                    <th data-field="prefix">prefix</th>
                    <th data-field="incoming">incoming</th>
                    <th data-field="outgoing">outgoing</th>
                </tr>
              </thead>
              <tbody>
              <tr v-for="redirekt in filterRedirekts(filter)">
              <td><a class="btn waves" v-modal:edit v-on:click="selectRedirekt(redirekt)">{{ redirekt.prefix }}</a></td>
              <td>{{ redirekt.incoming }}</td>
              <td>{{ redirekt.outgoing }}</td>
              </tr>
              </tbody>
            </table>
            </div>
          </div>


          <v-modal id="edit">
            <div slot="content">
              <h4>edit redirekt</h4>
              <table id="tableredirekts">
              <thead>
                <tr>
                    <th data-field="prefix">prefix</th>
                    <th data-field="incoming">incoming</th>
                    <th data-field="outgoing">outgoing</th>
                </tr>
              </thead>
              <editor v-bind:selection=this.selection ></editor>
              </table>
            </div>
            <div slot="footer">
              <v-btn-link class="pull-right red" waves-light modal flat>Delete</v-btn-link>
              <v-btn-link class="pull-right green" waves-light modal flat>Update</v-btn-link>
            </div>
          </v-modal>

      </div>
    </div> <!-- container -->

  </div>
</template>

<script>
export default {
  name: 'app',
  components: {
    'editor': {
      name: 'editor',
      props: ['selection'],
      template: '<tr><td>{{ selection.prefix }}</td><td>{{ selection.incoming }}</td><td>{{ selection.outgoing }}</td></tr>'
    }
  },
  data () {
    return {
      msg: 'making your redirect life easier.',
      redirektdata: null,
      selection: { prefix: '', incoming: '', outgoing: '' },
      loading: true,
      filter: '',
      error: null
    }
  },
  methods: {
    redirektsOk: function(response) {
          this.redirektdata = response.body.response;
          this.loading = false;
          console.log(this.redirektdata);
    },
    redirektsError: function(response) {
          this.loading = false;
          this.error = 'there was an error';
          console.log(JSON.stringify(response.statusText));
    },
    selectRedirekt: function(redir) {
          this.selection = redir;
          console.log(JSON.stringify(redir));
    },
    filterRedirekts: function(value) {
            return this.redirektdata.redirekts.filter(function(item) {
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
#tableredirekts {
  margin-left: 1em;
}
</style>

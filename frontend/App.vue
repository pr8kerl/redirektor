
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

      <div class="section no-pad-bot center" id="index-banner">
        <h2 class="header center orange-text text-darken-2">{{ msg }}</h2>

        <!--[if lt IE 8]>
        <div class="row center pink-text"><h2>You are using an <strong>outdated</strong> browser. Please <a href="http://browsehappy.com/">upgrade your browser</a> to improve your experience.</h2></div>
        <![endif]-->

          <div v-if="loading">
            <v-progress-circular active red red-flash></v-progress-circular>
          </div>

          <div v-if="redirektdata" class="container">
            <div class="row">
              <div class="col s1">
                <i class="material-icons prefix">search</i>
              </div>
              <div class="col s10">
                <input v-model="filter" class="form-control" placeholder="filter incoming">
              </div>
              <div class="col s1">
                <a v-modal:add v-on:click="initAddRedirekt()"><i class="small material-icons">add_circle_outline</i></a></td>
              </div>
            </div>
            <div class="row">
            <table id="tableredirekts" class="striped">
              <thead>
                <tr>
                    <th data-field="edit">edit</th>
                    <th data-field="prefix">prefix</th>
                    <th data-field="incoming">incoming</th>
                    <th data-field="outgoing">outgoing</th>
                </tr>
              </thead>
              <tbody>
              <tr v-for="redirekt in filterIncoming(filter)">
              <td><a v-modal:edit v-on:click="selectRedirekt(redirekt)"><i class="small material-icons">mode_edit</i></a></td>
              <td>{{ redirekt.prefix }}</td>
              <td>{{ redirekt.incoming }}</td>
              <td>{{ redirekt.outgoing }}</td>
              </tr>
              </tbody>
            </table>
            </div>
          </div>


          <div id="edit" class="modal">
            <div class="modal-content">
              <h4>edit redirekt</h4>
                <editor v-bind:selection=this.selection></editor>
            </div>
            <div class="modal-footer">
              <a class="modal-action modal-close waves-effect green btn-flat" v-on:click="putRedirekt(selection)" >Update</a>
              <a class="modal-action modal-close waves-effect red btn-flat" v-on:click="deleteRedirekt(selection)" >Delete</a>
            </div>
          </div>

          <div id="add" class="modal">
            <div class="modal-content">
              <h4>add redirekt</h4>
                <div id="addredirekt" class="container"> 
                  <div class="row left-align valign-wrapper"> 
                    <span class="col s2 theader">prefix</span>
                    <span class="col s10">
                      <div class="input-field">
                        <div class="input-field"> 
                          <v-select v-model="addition.prefix" name="select" id="select" :items="prefixes" ></v-select>
                        </div>
                      </div>
                    </span>
                  </div>
                  <div class="row left-align valign-wrapper">
                    <span class="col s2 theader">incoming</span><span class="col s10"><input v-model="addition.incoming" class="form-control" placeholder="incoming url"></span>
                  </div>
                  <div class="row left-align valign-wrapper">
                    <span class="col s2 theader">outgoing</span><span class="col s10"><input v-model="addition.outgoing" class="form-control" placeholder="outgoing url"></span>
                  </div>
                </div>
            </div>
            <div class="modal-footer">
              <a class="modal-action modal-close waves-effect green btn-flat" v-on:click="addRedirekt(addition)" >Add</a>
              <a class="modal-action modal-close waves-effect orange btn-flat" >Cancel</a>
            </div>
          </div>

      </div>

  </div>
</template>

<script>
export default {
  name: 'app',
  components: {
    'editor': {
      name: 'editor',
      props: ['selection' ],
      template: '<div id="editredirekt" class="container"> <div class="row left-align valign-wrapper"> <span class="col s2 left-align">prefix</span><span class="col s10">{{ selection.prefix }}</span> </div> <div class="row left-align valign-wrapper"> <span class="col s2">incoming</span><span class="col s10">{{ selection.incoming }}</span> </div> <div class="row left-align valign-wrapper"> <span class="col s2">outgoing</span><span class="col s10"><input v-model="selection.outgoing" class="form-control" placeholder="selection.outgoing"></span> </div></div>'
    }
  },
  data () {
    return {
      msg: 'making your redirect life easier.',
      redirektdata: null,
      selection: { prefix: '', incoming: '', outgoing: '' },
      addition: { prefix: '', incoming: '', outgoing: '' },
      loading: true,
      filter: '',
      status: null,
      prefixes: [ { id: 1, text: "aem", }, { id: 2, text: "aemau", }, { id: 3, text: "aemnz", }, { id: 4, text: "fw", }, { id: 5, text: "fwau", }, { id: 6, text: "fwnz" } ],
      error: null
    }
  },
  methods: {
    redirektsOk: function(response) {
          this.redirektdata = response.body.response.redirekts;
          this.selection = { prefix: '', incoming: '', outgoing: '' };
          this.loading = false;
          this.error = null;
          this.status = null;
    },
    redirektsError: function(response) {
          this.loading = false;
          this.error = response.body.error;
          Materialize.toast(this.error, 5000, 'rounded red');
          console.log(JSON.stringify(response.statusText));
    },
    initAddRedirekt: function() {
          this.status = '';
          this.error = null;
          console.log(JSON.stringify(this.prefixes));
    },
    addRedirekt: function(redir) {
        switch(redir.prefix) {
          case "1":
            redir.prefix="aem";
            break;
          case "2":
            redir.prefix="aemau";
            break;
          case "3":
            redir.prefix="aemnz";
            break;
          case "4":
            redir.prefix="fw";
            break;
          case "5":
            redir.prefix="fwau";
            break;
          case "6":
            redir.prefix="fwnz";
            break;
          default:
        }
        console.log(JSON.stringify(redir));
        this.$http.put('/api/redirekts', redir).then(this.addOk,this.redirektsError);
    },
    selectRedirekt: function(redir) {
          this.status = '';
          this.error = null;
          if (redir !== undefined) {
            this.selection = redir;
            console.log("selected: " + JSON.stringify(redir));
          } else {
            this.selection = { prefix: '', incoming: '', outgoing: '' };
          }
    },
    filterIncoming: function(value) {
            return this.redirektdata.filter(function(item) {
                return item.incoming.toLowerCase().indexOf(value) > -1;
            });
    },
    putOk: function(response) {
          this.status = "redirect updated successfully";
          Materialize.toast(this.status, 5000, 'rounded green');
          this.loading = true;
          this.selection = { prefix: '', incoming: '', outgoing: '' };
          this.redirektdata = null;
          this.error = null;
          this.filter = '';
          this.$http.get('/api/redirekts').then(this.redirektsOk,this.redirektsError);
    },
    addOk: function(response) {
          this.status = "redirect created successfully";
          Materialize.toast(this.status, 5000, 'rounded green');
          this.loading = true;
          this.addition = { prefix: '', incoming: '', outgoing: '' };
          this.redirektdata = null;
          this.error = null;
          this.filter = '';
          this.$http.get('/api/redirekts').then(this.redirektsOk,this.redirektsError);
    },
    putRedirekt: function(redir) {
        console.log(JSON.stringify(redir));
        this.$http.put('/api/redirekts', redir).then(this.putOk,this.redirektsError);
    },
    deleteOk: function(response) {
          this.status = "redirect deleted successfully";
          Materialize.toast(this.status, 5000, 'rounded green');
          this.loading = true;
          this.redirektdata = null;
          this.error = null;
          this.filter = '';
          this.$http.get('/api/redirekts').then(this.redirektsOk,this.redirektsError);
    },
    deleteRedirekt: function(redir) {
        console.log("delete: " + JSON.stringify(redir));
        this.$http.patch('/api/redirekts', redir).then(this.deleteOk,this.redirektsError);
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

a {
  color: #b71c1c;
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

.theader {
  font-weight: bold;
}

#tableredirekts {
}

/* label color */
.input-field label {
	color: #2c3e50;
}
/* label focus color */
.input-field input[type=text]:focus + label {
	color: #2c3e50;
}
/* label underline focus color */
.input-field input[type=text]:focus {
	border-bottom: 1px solid #2c3e50;
	box-shadow: 0 1px 0 0 #2c3e50;
}
/* valid color */
.input-field input[type=text].valid {
	border-bottom: 1px solid #2c3e50;
	box-shadow: 0 1px 0 0 #2c3e50;
}
/* invalid color */
.input-field input[type=text].invalid {
	border-bottom: 1px solid red;
	box-shadow: 0 1px 0 0 red;
}
/* icon prefix focus color */
	.input-field .prefix.active {
	color: #2c3e50;
}
</style>

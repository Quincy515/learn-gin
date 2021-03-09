<template>
  <div class="login" style="width:500px;margin:10px auto">
    <Form :model="formItem" :label-width="80">
      <FormItem label="用户名">
        <Input v-model="formItem.userName" placeholder="输入用户名"/>
      </FormItem>
      <FormItem label="密码">
        <Input v-model="formItem.userPass"  type="password" placeholder="密码"/>
      </FormItem>
      <FormItem>
        <Button type="primary">登录</Button>

      </FormItem>
    </Form>
  </div>
</template>

<script>
  import { Form,Input,FormItem,Button } from 'iview';
export default {
  name: 'LoginForm',
  data(){
    return {
      formItem:{
        userName:"",
        userPass:""
      }
    }
  },
  props: {

  },
  components:{
    Form,
    Input,
    FormItem,
    Button
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
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

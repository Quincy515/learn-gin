<template>
    <div>
      <div style="width: 95%;margin: 0 auto">
        <div style="width: 95%;margin:30px auto">
            <el-table
                    :data="podList"
                    border
                    style="width: 100%">
                <el-table-column
                        prop="PodName"
                        label="PodName"
                        width="100">
                </el-table-column>
                <el-table-column
                        prop="PodImage"
                        label="Pod镜像"
                        width="220">
                </el-table-column>
                <el-table-column
                        prop="PodNode"
                        label="Pod节点"
                        width="100">
                </el-table-column>

            </el-table>
        </div>

    </div>
    </div>
</template>
<script type="module">
   module.exports ={
       data(){
           return {
               podList: [],
           }
       },
       created(){
            NewClient().onmessage=(e)=>{
                if(e.data !== 'ping'){
                    this.podList = JSON.parse(e.data)
                    this.$forceUpdate()
                }
            }
       },


   }
</script>
<style>
  .sdt{margin: 10px auto;width:90%;border-radius: 5px;display: block;float:left;margin-left: 50px}
  .sdt dt{width:100%;display: block;color:#3A7B43;font-size:16px;font-weight: bold;margin-bottom: 10px}
  .sdt dd{float:left;margin:0 auto;text-indent:1em}
  .sdt .row{width:100%;}
  .sdt dd .search{width:50%;}
  .sdt dd a{color: #3a8ee6}
  .sdt dd a:hover{background: #eb5975;color:#fff}
  .sdt dd .select{background: #eb5975;color:#fff}
  .sdt dd .numtext{width:100px}
  .sdt dd span{margin: 0 auto}
    a{cursor:pointer}
    .el-pagination{margin:0 auto;float: left}
</style>
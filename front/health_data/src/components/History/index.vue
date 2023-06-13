<template>
    <el-dialog title="历史检查项目" :visible.sync="show" width="70%" modal :show-close="false" :close-on-press-escape="false" :close-on-click-modal="false">
        <div>
            <el-table
                :data="tableData"
                border
                :header-cell-style="$handerStyleEvent" 
                style="width: 100%"
                >
                <el-table-column
                    type="index"
                    align="center"
                    width="50">
                </el-table-column>
                <el-table-column 
                    v-for="(item,i) in column"
                    :key="i"
                    :prop="item.prop"
                    :label="item.label"
                    :align="item.align"
                >
                </el-table-column>
                <el-table-column fixed="right" label="操作" width="100" align="center">
                    <template slot-scope="scope">
                        <span v-if="false">{{ scope }}</span>
                        <el-button type="text" size="small" @click="viewClick">查看</el-button>
                    </template>
                </el-table-column>
            </el-table>
        </div>
        <span slot="footer" class="dialog-footer">
            <el-button type="primary" @click="onClick">关 闭</el-button>
        </span>
    </el-dialog>
</template>

<script>
export default {
    props:{
        show:{
            type:Boolean,
            default:false,
        },
        tableData:{
            type:Array,
            default:[],
        }
    },
    data(){
        return{
            column:[
                {
                    prop:'transition',
                    label:'transitionID',
                    align:'center'
                },{
                    prop:'',
                    label:'检查项目',
                    align:'center'
                },{
                    prop:'',
                    label:'是否正常',
                    align:'center'
                },
            ]
        }
    },
    methods:{      
        onClick(){
            this.$emit('close',{open:false});
        },
        viewClick(){
            
        },
    },
}
</script>
<template>
  <div class="home">
    <div class="heard">
      <div class="login">
        <div v-if="this.form.addr == null" class="Register">
          <span @click="Register">创建账户</span>
        </div>
        <span class="userName" v-else>{{ this.form.addr }}</span>
      </div>
    </div>
    <div class="main">
      <div class="content">
        <div class="left">
          <img src="../assets/humanBody.jpg" alt="加载失败...">
        </div>
        <div class="right">
          <div class="form">
            <div class="form_data">
              <div class="row">
                <span class="label">
                  检查项目
                </span>
                <el-select v-model="form.classify" placeholder="请选择" style="width:70%" :disabled="!form.addr">
                  <el-option v-for="item in options" :key="item.value" :label="item.label" :value="item.value">
                  </el-option>
                </el-select>
              </div>
              <div class="row" v-for="(item, i) in formDate" :key="i">
                <span class="label">
                  {{ item.label }}
                </span>
                <el-input class="input" v-if="form.addr" v-model="form[item.valueKey]" :placeholder="item.placeholder"></el-input>
                <el-input class="input" v-else v-model="form[item.valueKey]" placeholder="请先创建用户" disabled></el-input>
              </div>
            </div>
            <div class="operate">
              <button class="btn" @click="submit">提交</button>
              <button class="btn" @click="historyClick" v-if="tableData.length > 0">历史检查</button>
            </div>
          </div>
        </div>
      </div>
    </div>
    <history :tableData="tableData" :show="show" @close="close"></history>
  </div>
</template>
<script>
import modal from '@/plugins/modal'
import { Message } from 'element-ui';
import { createAccount, submitData } from '../api/home/index.js'
import history from '../components/History/index.vue'
export default {
  components: {
    history
  },
  name: 'HomeView',
  data() {
    return {
      form: {
        classify: 1,
        addr: null,
      },
      options: [
        {
          label: '测血压',
          value: 1,
        }, {
          label: '测心率',
          value: 2,
        }, {
          label: '测血糖',
          value: 3,
        },
      ],
      formDate: [
        {
          label: '收缩压',
          valueKey: 'dbp',
          placeholder: '请输入收缩压',
        }, {
          label: '舒张压',
          valueKey: 'sbp',
          placeholder: '请输入舒张压',
        }, {
          label: '心率',
          valueKey: 'hr',
          placeholder: '请输入心率',
        }, {
          label: '血糖',
          valueKey: 'pbg',
          placeholder: '请输入血糖',
        },
      ],
      tableData: [{}],
      show: false,
    }
  },
  methods: {
    historyClick() {
      this.show = true
    },
    close(val) {
      this.show = val.open
    },
    Register() {
      modal.loading("创建用户中...");
      createAccount().then(res => {
        if (res.code == 0) {
          this.form.addr = res.data.address
        }
        modal.closeLoading();
      })
    },
    submit() {
      if (this.form.addr == null || this.form.addr == undefined || this.form.addr == "") {
        Message.error("请先创建账户")
      } else {
        modal.loading("正在请求...");
        for (let vattr in this.form) {
          if (vattr != 'addr') {
            this.form[vattr] = Number(this.form[vattr])
          }
        }
        submitData(this.form).then(res => {
          console.log(res);
          // if(res.code == 0){
          // }
          modal.closeLoading();
        })
      }
    },
  },
};
</script>
<style lang="less" scoped>
.home {
  .heard {
    background-color: cadetblue;
    width: 100%;
    height: 5vh;
    display: flex;
    align-items: center;
    .login {
      width: 70%;
      margin: 0 auto;
      display: flex;
      justify-content: flex-end;
      color: #fff;
      .Register {
        cursor: pointer;

        &:hover {
          color: #2c7bd4;
        }
      }

      .userName {
        cursor: pointer;
      }
    }
  }

  .main {
    width: 100%;
    height: 95vh;
    background-color: #FBF6F0;

    .content {
      width: 70%;
      margin: 0 auto;
      display: flex;

      .left {
        width: 40%;
        height: 95vh;

        img {
          width: 100%;
          height: 90vh;
        }
      }

      .right {
        width: 60%;
        display: flex;
        align-items: center;

        .form {
          width: 55%;
          margin: 0 auto;
          background-color: #fff;
          border-radius: 10px;
          box-shadow: 3px 3px 10px 10px #bfbfbf;

          .form_data {
            width: 100%;

            .row {
              width: 100%;
              height: 5vh;
              display: flex;
              justify-content: space-around;
              align-items: center;
              padding-top: 3vh;

              .label {
                width: 20%;
                font-size: 20px;
                font-weight: 700;
                text-align: justify;
                text-align-last: justify;
              }

              .input,
              .el-input {
                width: 70%;
                height: 36px;
              }
            }
          }

          .operate {
            display: flex;
            justify-content: space-around;
            align-items: center;
            padding: 2vh 0vh 2vh 0vh;

            .btn {
              width: 150px;
              height: 45px;
              border: none;
              background-color: #409EFF;
              color: #fff;
              border-radius: 10px;
              cursor: pointer;

              &:hover {
                background-color: #66B1FF;
              }
            }
          }
        }
      }
    }
  }
}</style>
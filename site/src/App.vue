<template>
  <div id="app">
    <a class="fork" href="https://github.com/vicanso/tiny">
      <img src="./assets/logo.png">
      <span>fork me</span>
    </a>
    <el-upload
      class="upload-demo"
      drag
      action="/api/upload"
      :on-error="handleError"
      :on-success="handleSuccess"
      :on-remove="handleRemove"
      :on-exceed="handleExceed"
      :limit="1"
    >
      <i class="el-icon-upload"></i>
      <div class="el-upload__text">Drop file here or <em>click to upload</em></div>
      <div class="el-upload__tip" slot="tip">files with a size less than 1mb (only support one file)</div>
    </el-upload>
    <div class="tinyMode">
      <el-radio
        v-for="(item, index) in modeList"
        v-model="mode"
        :label="index"
        :key="item"
      >
        {{item}}
      </el-radio>
    </div>
    <div class="tinyQuality">
      <el-input
        placeholder="Please input quality"
        v-model="quality"
      >
        <template slot="prepend">Quality:</template>
      </el-input>
      <div class="el-upload__tip" slot="tip">Webp lossless type should set quality to 0</div>
    </div>
    <div
      class="tinyResult"
      v-if="result"
    >
      orignal size: {{result.originalSize.toLocaleString()}}
      new size: {{result.size.toLocaleString()}}
      <el-progress :percentage="result.percent"></el-progress>
    </div>
    <el-button
      type="primary"
      :disabled="!file || mode === -1"
      plain
      @click="doTiny"
    >Do Tiny</el-button>
    <el-button
      type="success"
      :disabled="!result"
      plain
      @click="doDownload"
    >Download</el-button>
  </div>
</template>

<script>
import request from 'axios'
import {
  Loading
} from 'element-ui'
export default {
  name: 'app',
  data() {
    return {
      mode: -1,
      quality: null,
      file: '',
      status: '',
      result: null,
      modeList: [
        'gzip',
        'brotli',
        'jpeg',
        'png',
        'webp',
        'guetzli',
      ],
    };
  },
  watch: {
    mode(v) {
      switch (v) {
        case 0:
        case 1:
          this.quality = 9;
          break;
        case 2:
        case 5:
          this.quality = 90;
          break;
        case 3:
          this.quality = 80;
          break;
        case 4:
          this.quality = 75;
          break;
        default:
          this.quality = null;
          break;
      }
      this.reset();
    },
    quality() {
      this.reset();
    }
  },
  methods: {
    handleSuccess(res) {
      this.file = res.file;
    },
    reset() {
      this.result = null;
    },
    async doTiny() {
      if (this.status === 'doing') {
        return;
      }
      this.status = 'doing';
      const loadingInstance = Loading.service();
      try {
        const res = await request.post('/api/tiny', {
          file: this.file,
          mode: this.mode,
          quality: parseInt(this.quality, 10),
        });
        const {
          data,
        } = res;
        data.percent = Math.ceil(data.size * 100 / data.originalSize);
        this.result = data; 
      } catch (err) {
        let message = err.message;
        if (err.response && err.response.data && err.response.data.message) {
          message = err.response.data.message;
        }
        this.$notify({
          title: 'Error',
          message: `do tiny fail, ${message}`,
        });
      } finally {
        this.status = '';
        loadingInstance.close();
      }
    },
    doDownload() {
      if (!this.result) {
        return;
      }
      const {
        file,
      } = this.result;
      window.location.href = `/api/download/${file}`;
    },
    handleExceed() {
      this.$notify({
        title: 'Warning',
        message: 'You should remove the current file first',
      });
    },
    handleError(err) {
      this.$notify({
        title: 'Error',
        message: err.message,
      });
    },
    handleRemove() {
      this.reset();
      this.file = '';
    },
  },
}
</script>

<style>
@font-face {
  font-family: "Oleo Script";
  src: url("https://fonts.googleapis.com/css?family=Oleo+Script:700");
}
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
.fork {
  display: block;
  width: 200px;
  margin: 60px auto 20px auto;
  text-decoration: none;
}
.fork img {
  display: block;
  margin: auto;
}
.fork span {
  font-family: Oleo Script,cursive;
  text-transform: uppercase;
  font-size: 24px;
  line-height: 2;
  color: #0f0f0f;
}
.tinyMode {
  margin: 20px 0;
}
.tinyQuality {
  margin: 20px auto;
  width: 500px;
}
.tinyResult {
  text-align: center;
  margin: 20px auto;
  width: 500px;
}
</style>

<template>
  <div class="data-screen">
    <!-- Full Screen Container for DataV -->
    <dv-full-screen-container>
      <!-- Header -->
      <div class="screen-header">
        <dv-decoration-8 style="width:300px;height:50px;" />
        <div class="header-center">
          <dv-decoration-5 style="width:500px;height:40px;" />
          <div class="title-text">智慧社区数据可视化大屏</div>
        </div>
        <dv-decoration-8 :reverse="true" style="width:300px;height:50px;" />
        <div class="time-text">{{ currentTime }}</div>
        <div class="back-btn" @click="$router.push('/home')">
            <el-icon><HomeFilled /></el-icon> 首页
        </div>
      </div>

      <!-- Body -->
      <div class="screen-body">
        <el-row :gutter="20" style="height: 100%">
          <!-- Left Column -->
          <el-col :span="6" class="column-side">
            <!-- Card 1: Core Metrics -->
            <dv-border-box-11 class="tech-card h-30" title="核心指标">
              <div class="metrics-grid p-20">
                <div class="metric-box">
                  <div class="metric-label">总用户</div>
                  <div class="metric-value">{{ stats.totalUsers || 0 }}</div>
                </div>
                <div class="metric-box">
                  <div class="metric-label">日订单</div>
                  <div class="metric-value">{{ stats.todayOrders || 0 }}</div>
                </div>
                <div class="metric-box">
                  <div class="metric-label">车位率</div>
                  <div class="metric-value">{{ stats.parkingRate || '0%' }}</div>
                </div>
                <div class="metric-box">
                  <div class="metric-label">月营收</div>
                  <div class="metric-value small-font">¥{{ stats.monthIncome || 0 }}</div>
                </div>
              </div>
            </dv-border-box-11>

            <!-- Card 2: Repair Stats -->
            <dv-border-box-13 class="tech-card h-35 mt-15">
               <div class="chart-title">报修类型占比</div>
               <div ref="pieChartRef" class="chart-container"></div>
            </dv-border-box-13>

            <!-- Card 3: Energy Trend -->
            <dv-border-box-13 class="tech-card h-35 mt-15">
               <div class="chart-title">水电能耗趋势</div>
               <div ref="lineChartRef" class="chart-container"></div>
            </dv-border-box-13>
          </el-col>

          <!-- Center Column -->
          <el-col :span="12" class="column-center">
             <div class="center-map-box">
               <dv-border-box-10>
                   <!-- Background Image -->
                   <div class="map-bg"></div>
                   <div class="map-title-overlay">社区全景模型</div>
                   
                   <!-- Center Stats Overlay -->
                   <div class="center-data">
                      <div class="c-item">
                         <div class="c-label">年度总交易额</div>
                         <div class="c-num"><dv-digital-flop :config="flopConfig1" style="width:200px;height:50px;" /></div>
                      </div>
                       <div class="c-item">
                         <div class="c-label">安保巡逻次数</div>
                         <div class="c-num"><dv-digital-flop :config="flopConfig2" style="width:200px;height:50px;" /></div>
                      </div>
                   </div>
               </dv-border-box-10>
             </div>
          </el-col>

          <!-- Right Column -->
          <el-col :span="6" class="column-side">
             <!-- Card 1: Cost Analysis -->
             <dv-border-box-13 class="tech-card h-40">
               <div class="chart-title">费用构成分析</div>
               <div ref="barChartRef" class="chart-container"></div>
             </dv-border-box-13>

             <!-- Card 2: Visitor Monitor -->
             <dv-border-box-11 class="tech-card h-60 mt-15" title="实时访客监控">
               <div class="visitor-list p-20">
                   <dv-scroll-board :config="scrollBoardConfig" style="width:100%;height:100%" />
               </div>
             </dv-border-box-11>
          </el-col>
        </el-row>
      </div>
    </dv-full-screen-container>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, reactive, watch } from 'vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import { getDashboardStats } from '@/api/admin'
import { HomeFilled } from '@element-plus/icons-vue'

const currentTime = ref(dayjs().format('YYYY-MM-DD HH:mm:ss'))
let timer = null
const stats = ref({})

// DataV Configs
const flopConfig1 = reactive({
    number: [0],
    content: '¥ {nt}',
    style: { fontSize: 30, fill: '#ffd04b' }
})
const flopConfig2 = reactive({
    number: [0],
    style: { fontSize: 30, fill: '#ffd04b' }
})
const scrollBoardConfig = reactive({
    header: ['姓名', '时间', '状态'],
    data: [],
    rowNum: 7,
    headerBGC: 'rgba(0, 242, 254, 0.1)',
    oddRowBGC: 'rgba(0, 0, 0, 0.1)',
    evenRowBGC: 'rgba(0, 242, 254, 0.05)',
    columnWidth: [80, 150, 80],
    align: ['center', 'center', 'center']
})

const pieChartRef = ref(null)
const lineChartRef = ref(null)
const barChartRef = ref(null)

let pieChart = null
let lineChart = null
let barChart = null

const fetchData = async () => {
    try {
        const res = await getDashboardStats()
        if (res) {
             stats.value = res
             // Update Flops
             flopConfig1.number = [res.yearTotalAmount || 0]
             flopConfig2.number = [res.patrolCount || 0]

             // Update Scroll Board
             if (res.visitorLogs && res.visitorLogs.length > 0) {
                 scrollBoardConfig.data = res.visitorLogs.map(item => [
                     item.name || '访客',
                     item.visit_time ? item.visit_time.slice(5, 16).replace('T', ' ') : '--',
                     item.status === 1 ? '<span style="color:#67c23a">通过</span>' : '<span style="color:#e6a23c">待审</span>'
                 ])
             } else {
                 scrollBoardConfig.data = [['暂无', '--', '--']]
             }

             initCharts(res)
        }
    } catch (e) {
        console.error("Dashboard data fetch failed", e)
    }
}

const initCharts = (data) => {
    // 1. Pie Chart
    if(pieChart) pieChart.dispose()
    pieChart = echarts.init(pieChartRef.value)
    const pieData = data.repairStats && data.repairStats.length > 0 ? data.repairStats : [{name: '无数据', value: 0}]
    pieChart.setOption({
        color: ['#00f2fe', '#409EFF', '#e6a23c', '#f56c6c'],
        legend: { bottom: 0, textStyle: { color: '#fff' } },
        series: [{
            type: 'pie',
            radius: ['40%', '60%'],
            center: ['50%', '45%'],
            label: { color: '#fff' },
            data: pieData
        }]
    })

    // 2. Line Chart
    if(lineChart) lineChart.dispose()
    lineChart = echarts.init(lineChartRef.value)
    lineChart.setOption({
        tooltip: { trigger: 'axis' },
        grid: { top: '15%', left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: {
            type: 'category',
            boundaryGap: false,
            data: data.incomeDates,
            axisLabel: { color: '#fff' }
        },
        yAxis: {
            type: 'value',
            splitLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } },
            axisLabel: { color: '#fff' }
        },
        series: [{
            type: 'line',
            smooth: true,
            lineStyle: { width: 2, color: '#00f2fe' },
            areaStyle: { color: 'rgba(0, 242, 254, 0.3)' },
            data: data.incomeTrend
        }]
    })

    // 3. Bar Chart
    if(barChart) barChart.dispose()
    barChart = echarts.init(barChartRef.value)
    barChart.setOption({
        tooltip: { trigger: 'axis' },
        grid: { top: '20%', left: '3%', right: '4%', bottom: '3%', containLabel: true },
        xAxis: {
            type: 'category',
            data: ['物业', '停车', '商城'],
            axisLabel: { color: '#fff' }
        },
        yAxis: {
             type: 'value',
             splitLine: { lineStyle: { color: 'rgba(255,255,255,0.1)' } },
             axisLabel: { color: '#fff' }
        },
        series: [{
            type: 'bar',
            barWidth: '30%',
            itemStyle: { color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [{ offset: 0, color: '#00f2fe' }, { offset: 1, color: '#409EFF' }]) },
            data: data.costStructure || [0, 0, 0]
        }]
    })
}

const handleResize = () => {
    pieChart && pieChart.resize()
    lineChart && lineChart.resize()
    barChart && barChart.resize()
}

onMounted(() => {
    timer = setInterval(() => {
        currentTime.value = dayjs().format('YYYY-MM-DD HH:mm:ss')
    }, 1000)
    fetchData()
    window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
    clearInterval(timer)
    window.removeEventListener('resize', handleResize)
    pieChart && pieChart.dispose()
    lineChart && lineChart.dispose()
    barChart && barChart.dispose()
})
</script>

<style scoped>
.data-screen {
    width: 100vw;
    height: 100vh;
    background: #050a15;
    color: #fff;
    overflow: hidden;
    font-family: "Microsoft YaHei", sans-serif;
}

.screen-header {
    height: 60px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    position: relative;
    z-index: 10;
}
.header-center {
    display: flex;
    flex-direction: column;
    align-items: center;
    position: absolute;
    left: 50%;
    transform: translateX(-50%);
    top: 5px;
}
.title-text {
    font-size: 28px;
    font-weight: bold;
    color: #00f2fe;
    margin-top: -30px;
    text-shadow: 0 0 10px #00f2fe;
}
.time-text {
    position: absolute;
    right: 320px;
    top: 15px;
    color: #fff;
    font-weight: bold;
}
.back-btn {
    position: absolute;
    left: 320px;
    top: 15px;
    color: #00f2fe;
    cursor: pointer;
    font-weight: bold;
    display: flex;
    align-items: center;
    gap: 5px;
}

.screen-body {
    height: calc(100vh - 60px);
    padding: 10px 20px 20px 20px;
    box-sizing: border-box;
}

.column-side, .column-center {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.tech-card {
    position: relative;
    box-sizing: border-box;
}
.h-30 { height: 30%; }
.h-35 { height: 35%; }
.h-40 { height: 40%; }
.h-60 { height: 60%; }
.mt-15 { margin-top: 15px; }
.p-20 { padding: 40px 20px 10px 20px; height: 100%; box-sizing: border-box; }

.chart-title {
    position: absolute;
    top: 10px;
    left: 20px;
    font-size: 16px;
    font-weight: bold;
    color: #fff;
    padding-left: 10px;
    border-left: 3px solid #00f2fe;
}
.chart-container {
    width: 100%;
    height: 100%;
    padding-top: 30px;
    box-sizing: border-box;
}

.metrics-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 10px;
}
.metric-box {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: rgba(0, 242, 254, 0.05);
}
.metric-label { font-size: 14px; color: #ccc; }
.metric-value { font-size: 24px; color: #00f2fe; font-weight: bold; }

.center-map-box {
    width: 100%;
    height: 100%;
    position: relative;
    padding: 2px;
}
.map-bg {
    width: 100%;
    height: 100%;
    background-image: url('../../assets/images/cyber_city.png');
    background-size: cover;
    background-position: center;
    opacity: 0.9;
}
.map-title-overlay {
    position: absolute;
    top: 30px;
    width: 100%;
    text-align: center;
    font-size: 24px;
    color: #fff;
    text-shadow: 0 0 10px #00f2fe;
    z-index: 2;
}
.center-data {
    position: absolute;
    bottom: 60px;
    width: 100%;
    display: flex;
    justify-content: space-around;
    z-index: 2;
}
.c-label { font-size: 16px; color: #ccc; text-align: center; margin-bottom: 5px; }

</style>

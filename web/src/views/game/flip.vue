<template>
  <div class="game-flip">
    <iframe :src="gameUrl" width="380" height="610" />

    <div class="gf-score">
      <el-collapse accordion>
        <div>Score Count:&nbsp;{{ gameScoreStore.count }}</div>
        <el-collapse-item v-for="(item, index) in gameScoreStore.scores" :key="index">
          <template #title>{{ index + 1 }}{{ ". " + item.player_name +" "+ item.score }}</template>

          <div>{{ displayFlipResult(item.result) }}</div>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue"
import { useGameScoreStore } from "@/pinia/game_score.ts"
import { GameName } from "@/axios/ts/game.go.ts"
import { isFlipResult } from "@/ts/data.ts"

const gameScoreStore = useGameScoreStore()

const gameUrl = ref<string>(import.meta.env.Vite_axios_flip_game_url)

onMounted(() => {
    gameScoreStore.listGameScore(GameName.Flip, 10, 1)
})

function displayFlipResult(flipResult: string): string {
    try {
        const obj = JSON.parse(flipResult)

        return isFlipResult(obj) ? `Duration: ${obj.duration}, Step: ${obj.steps}.` : flipResult
    } catch(error) {
        console.error("> Json Parse Failed. ", flipResult, error)
        return flipResult
    }
}
</script>

<style lang="less" scoped>
.game-flip {
	display: flex;

	.gf-score {
		margin-left: 8vw;
		width: 20vw;

		font-size: 1.6rem;
	}

	iframe {
		transform-origin: left top;
		@media (min-height: 1440px) {
			transform: scale(1.4, 1.4); // 大分辨率屏幕放大一些，1920*1080不用缩放
		}
	}
}
</style>

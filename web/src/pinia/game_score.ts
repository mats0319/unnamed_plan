import { defineStore } from "pinia"
import { GameName, GameScore, ListGameScoreRes } from "@/axios/ts/game.go.ts"
import { gameAxios } from "@/axios/ts/game.http.ts"
import { log } from "@/ts/log.ts"
import { ref } from "vue"
import { Pagination } from "@/axios/ts/common.go.ts"

export const useGameScoreStore = defineStore("game_score", () => {
    const currentGameName = ref<GameName>(GameName.placeholder)
    const count = ref<number>(0)
    const scores = ref<Array<GameScore>>(new Array<GameScore>())

    async function uploadGameScore(game_name: GameName, score: number, result: string, player: string): Promise<void> {
        await gameAxios.uploadGameScore(game_name, score, result, player)

        await listGameScore(game_name, 10, 1)

        log.success("Upload Game Score")
    }

    async function listGameScore(game_name: GameName, pageSize: number, pageNum: number): Promise<void> {
        const pagination: Pagination = { size: pageSize, num: pageNum }
        const { data }: { data: ListGameScoreRes } = await gameAxios.listGameScore(game_name, pagination)

        currentGameName.value = game_name
        count.value = data.count
        scores.value = data.scores

        log.success("List Game Score")
    }

    return { count, scores, uploadGameScore, listGameScore }
})

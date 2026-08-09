// Generate File, Should Not Edit.
// Author : mario. github.com/mats0319
// Code   : github.com/mats0319/study/go/gocts
// Version: gocts v0.2.5

import { axiosWrapper } from "./config"
import { AxiosResponse } from "axios"
import { UploadGameScoreRes, UploadGameScoreReq, GameName, ListGameScoreRes, ListGameScoreReq } from "./game.go"
import { Pagination } from "./common.go"

class GameAxios {
    public uploadGameScore(game_name: GameName, score: number, result: string, player: string): Promise<AxiosResponse<UploadGameScoreRes>> {
        const req: UploadGameScoreReq = {
            game_name: game_name,
            score: score,
            result: result,
            player: player,
        }

        return axiosWrapper.post("/game-score/upload", req)
    }

    public listGameScore(game_name: GameName, page: Pagination): Promise<AxiosResponse<ListGameScoreRes>> {
        const req: ListGameScoreReq = {
            game_name: game_name,
            page: page,
        }

        return axiosWrapper.post("/game-score/list", req)
    }
}

export const gameAxios: GameAxios = new GameAxios()

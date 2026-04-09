export interface CreateGameRequest {
   leagueId: string;
   homeTeamId: string;
   awayTeamId: string;
   date: string; // ISO8601
}

export type Half = "top" | "bottom";

export const isHalf = (v: unknown): v is Half => v === "top" || v === "bottom";

export interface GameStateDTO {
   inning: number;
   half: Half;
   outs: number;
   homeScore: number;
   awayScore: number;
   // optional: inningRuns may be present on some responses
   inningRuns?: number[];
}

export interface GameResponse {
   gameId: string;
   leagueId: string;
   homeTeamId: string;
   awayTeamId: string;
   rulesVersion: number;
   state: GameStateDTO;
}

export interface RecordAtBatRequest {
   gameId: string;
   leagueId?: string;
   playerId: string;
   teamId?: string;

   strikes?: number;
   balls?: number;
   fouls?: number;
   result?: string;
   rbi?: number;
}

export interface AtBatResponse {
   gameId: string;
   playerId: string;
   teamId: string;
   seq: number;

   inning: number;
   half: Half;
   strikes: number;
   balls: number;
   fouls: number;
   result: string;
   rbi: number;
}

export interface LeagueResponse {
   leagueId: string;
   name: string;
   currentRulesVersion: number;
}

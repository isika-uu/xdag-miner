export namespace models {
	
	export class APIConfig {
	    id?: string;
	    "worker-id"?: string;
	
	    static createFrom(source: any = {}) {
	        return new APIConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this["worker-id"] = source["worker-id"];
	    }
	}
	export class CPUConfig {
	    enabled: boolean;
	    "huge-pages": boolean;
	    "max-threads-hint": number;
	    priority?: number;
	    asm: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CPUConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this["huge-pages"] = source["huge-pages"];
	        this["max-threads-hint"] = source["max-threads-hint"];
	        this.priority = source["priority"];
	        this.asm = source["asm"];
	    }
	}
	export class HTTPConfig {
	    enabled: boolean;
	    host: string;
	    port: number;
	    "access-token"?: string;
	    restricted: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HTTPConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.host = source["host"];
	        this.port = source["port"];
	        this["access-token"] = source["access-token"];
	        this.restricted = source["restricted"];
	    }
	}
	export class MinerStatus {
	    running: boolean;
	    hashrate: number;
	    threads: number;
	    uptime: number;
	    pool: string;
	    algorithm: string;
	    connected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new MinerStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.hashrate = source["hashrate"];
	        this.threads = source["threads"];
	        this.uptime = source["uptime"];
	        this.pool = source["pool"];
	        this.algorithm = source["algorithm"];
	        this.connected = source["connected"];
	    }
	}
	export class PoolConfig {
	    algo?: string;
	    coin?: string;
	    url: string;
	    user: string;
	    pass: string;
	    "rig-id"?: string;
	    nicehash: boolean;
	    enabled: boolean;
	    tls: boolean;
	
	    static createFrom(source: any = {}) {
	        return new PoolConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.algo = source["algo"];
	        this.coin = source["coin"];
	        this.url = source["url"];
	        this.user = source["user"];
	        this.pass = source["pass"];
	        this["rig-id"] = source["rig-id"];
	        this.nicehash = source["nicehash"];
	        this.enabled = source["enabled"];
	        this.tls = source["tls"];
	    }
	}
	export class RandomXConfig {
	    init: number;
	    "init-avx2": number;
	    mode: string;
	    "1gb-pages": boolean;
	    numa: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RandomXConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.init = source["init"];
	        this["init-avx2"] = source["init-avx2"];
	        this.mode = source["mode"];
	        this["1gb-pages"] = source["1gb-pages"];
	        this.numa = source["numa"];
	    }
	}
	export class SystemInfo {
	    os: string;
	    arch: string;
	    cpuModel: string;
	    cpuCores: number;
	    totalMemory: number;
	    xmrigVersion: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.cpuModel = source["cpuModel"];
	        this.cpuCores = source["cpuCores"];
	        this.totalMemory = source["totalMemory"];
	        this.xmrigVersion = source["xmrigVersion"];
	    }
	}
	export class XMRigConfig {
	    api: APIConfig;
	    http: HTTPConfig;
	    autosave: boolean;
	    cpu: CPUConfig;
	    pools: PoolConfig[];
	    randomx: RandomXConfig;
	    "log-file"?: string;
	
	    static createFrom(source: any = {}) {
	        return new XMRigConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.api = this.convertValues(source["api"], APIConfig);
	        this.http = this.convertValues(source["http"], HTTPConfig);
	        this.autosave = source["autosave"];
	        this.cpu = this.convertValues(source["cpu"], CPUConfig);
	        this.pools = this.convertValues(source["pools"], PoolConfig);
	        this.randomx = this.convertValues(source["randomx"], RandomXConfig);
	        this["log-file"] = source["log-file"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}


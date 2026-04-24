export namespace dorks {
	
	export class Dork {
	    id: string;
	    name: string;
	    template: string;
	    category: string;
	    description: string;
	    severity: string;
	    needsDomain: boolean;
	    needsKeyword: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Dork(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.template = source["template"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.severity = source["severity"];
	        this.needsDomain = source["needsDomain"];
	        this.needsKeyword = source["needsKeyword"];
	    }
	}

}

export namespace main {
	
	export class DorkResult {
	    dork: dorks.Dork;
	    query: string;
	    results: search.SearchResult[];
	
	    static createFrom(source: any = {}) {
	        return new DorkResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dork = this.convertValues(source["dork"], dorks.Dork);
	        this.query = source["query"];
	        this.results = this.convertValues(source["results"], search.SearchResult);
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
	export class CategoryResult {
	    category: string;
	    dorks: DorkResult[];
	
	    static createFrom(source: any = {}) {
	        return new CategoryResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.category = source["category"];
	        this.dorks = this.convertValues(source["dorks"], DorkResult);
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
	
	export class SearchResponse {
	    mode: string;
	    modeLabel: string;
	    isMockMode: boolean;
	    isLiveMode: boolean;
	    fetchedCount: number;
	    categories: CategoryResult[];
	
	    static createFrom(source: any = {}) {
	        return new SearchResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mode = source["mode"];
	        this.modeLabel = source["modeLabel"];
	        this.isMockMode = source["isMockMode"];
	        this.isLiveMode = source["isLiveMode"];
	        this.fetchedCount = source["fetchedCount"];
	        this.categories = this.convertValues(source["categories"], CategoryResult);
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

export namespace search {
	
	export class SearchResult {
	    title: string;
	    url: string;
	    snippet: string;
	    resultType?: string;
	
	    static createFrom(source: any = {}) {
	        return new SearchResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.title = source["title"];
	        this.url = source["url"];
	        this.snippet = source["snippet"];
	        this.resultType = source["resultType"];
	    }
	}

}


// The 'health' program.
program health.aleo {
    struct HealthInfo {
    	classify: u32,
        sbp: u32,
        dbp: u32,
        hr: u32,
        pbg: u32,
        addr: address
    }

    record Health {
        owner: address,
        gates: u64,
        id: field,
        classify: u32,
        sbp: u32,
        dbp: u32,
        hr: u32,
        pbg: u32,
    }
    
    record HealthResult {
        owner: address,
        gates: u64,
        id: field,
        classify: u32,
        result: u32
    }

    //mapping healths: field => u64;

    transition submit(public info: HealthInfo) -> Health {       
        let id: field = BHP256::hash_to_field(info.addr);
        let result: u64 = 0u64;
        //if(info.classify == 1u32){
        //	result = measure_bp(info.sbp,info.dbp);
        //}else if(info.classify == 2u32){
        //	result = measure_hr(info.hr);
        //}else {
        //	result = measure_pbg(info.pbg);
        //}
        return Health {
            owner: self.caller,
            gates: 0u64,
            id,
            classify: info.classify,
            sbp: info.sbp,
            dbp: info.dbp,
            hr: info.hr,
            pbg: info.pbg,
        }; //then finalize(id,result);
    }
    //finalize submit(public id: field, result: u64) {
    //    Mapping::set(healths, id, result);
    //}
    
    //测血压
    transition measure_bp(public health: Health) -> HealthResult {
    	let result: u32 = 0u32;
    	let b: u32 = health.sbp;
    	let c: u32 = health.dbp;
        if((b >= 140u32 && b <= 159u32) || (c >=90u32 && c <= 99u32)){
            result = 1u32;
        }else if((b >= 160u32 && b <= 179u32) || (c >=100u32 && c <= 109u32)){
            result =  2u32;
        }else if(b >= 180u32 || c >=110u32){
            result =  3u32;
        }else {
            result =  0u32;
        }
        
        return HealthResult {
         owner: health.owner,
         gates: 0u64,
         id: health.id,
         classify: health.classify,
         result: result,
       };
    }
	
//测心率
    transition measure_hr(public health: Health) -> HealthResult {
    	let result: u32 = 0u32;
    	let b: u32 = health.hr;
        if(b >= 60u32 && b <= 100u32){
            result = 0u32;
        }else {
            result = 1u32;
        }
        
        return HealthResult {
         owner: health.owner,
         gates: 0u64,
         id: health.id,
         classify: health.classify,
         result: result,
       };
    }
    
    //测血糖
    transition measure_pbg(public health: Health) -> HealthResult {
    	let result: u32 = 0u32;
    	let b: u32 = health.pbg;
        if(b >= 39u32 && b <= 61u32){
            result = 0u32;
        }else {
            result = 1u32;
        }
        
        return HealthResult {
         owner: health.owner,
         gates: 0u64,
         id: health.id,
         classify: health.classify,
         result: result,
       };
    }
}

//
//  SignedMsg.swift
//  PDUM
//
//  Created by Liu Peng on 2021/8/31.
//

import Foundation

struct SignedMsg: Hashable, Codable, Identifiable {
    var content: String
    var refs: [String]
    var signature: String
    var id: String {
        signature
    }
    
    var photon: Photon? {
        let data = Data(base64Encoded: content)
        if let resObject = (try? JSONSerialization.jsonObject(with: data!)) as? [String: Any]{
            let t = resObject["t"] as! Int
            let v = resObject["v"] as! Int
            let d = resObject["d"] as! String
            return Photon(t:t,v:v,d:d)
        }
        return nil
    }
    
    var pInfo: PInfo? {
        let data = Data(base64Encoded: photon!.d)
        if let resObject = (try? JSONSerialization.jsonObject(with: data!)) as? [String: Any]{
            let text = resObject["text"] as! String
            let quote =  resObject["quote"] is NSNull ? nil : resObject["quote"] as? String
            let resources = resObject["resources"] as! [[String:Any]]
            var resArray = [PIRes]()
            resources.forEach { item in
                let ft = item["format"] as! Int
                let d = item["data"] is NSNull ? nil : Data(base64Encoded:item["data"] as! String)
                let url = item["url"] as! String
                let cs = item["cs"] as! String
                let pir = PIRes(format: ft, data: d, url: url, cs: cs)
                resArray.append(pir)
            }
            return PInfo(text: text, quote: quote, resources: resArray)
        }
        return nil
    }
    
    var pBorn: PBorn? = nil
    var pProfile: PProfile? = nil
    
    struct Photon : Hashable, Codable {
        var t: Int
        var v: Int
        var d: String
    }
    
    struct PIRes : Hashable, Codable {
        var format: Int
        var data: Data?
        var url: String
        var cs: String
    }
    
    struct PInfo : Hashable, Codable {
        var text: String
        var quote: String?
        var resources: [PIRes]
    }
    
    struct PBorn : Hashable, Codable {
        var addr: String
        var sigs: [String]
    }
    
    struct PProfile : Hashable, Codable {
        var name: String
        var email: String
        var bio: String
        var url: String
        var location: String
        var avatar: PIRes
        var extra: String
    }
}

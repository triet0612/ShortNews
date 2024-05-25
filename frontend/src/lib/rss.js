import {api_url} from "./index.js"
import { get } from "svelte/store"

export class NewsSource {
  /**
   * @param {string} pubID 
   * @param {string} link 
   * @param {string} voiceType 
   */
  constructor(pubID, link, voiceType) {
    this.pubID = pubID
    this.link = link
    this.voiceType = voiceType
  }
}
/**
 * @returns {Promise<NewsSource[]>}
 */
export async function newsSourcefromURL() {
  /**
   * @type {Array<any>}
   */
  let jsResponse = await fetch(get(api_url) + "/rss")
    .then(res => res.json())
    .catch(err => {console.log(err);return []})
  /**
   * @type {NewsSource[]}
   */
  let options = []
  jsResponse.map(v => {
    options.push(new NewsSource(
      v["PublisherID"], v["Link"], v["VoiceType"],
    ))
  })
  return options
}

/**
 * @param {NewsSource} newSrc 
 */
export async function createSource(newSrc) {
  let stat = await fetch(get(api_url) + "/rss", { 
    "method": "POST", 
    "body": JSON.stringify({
      "link": newSrc.link,
      "VoiceType": newSrc.voiceType
    })
  }).then(res => res.status === 200? "ok": "no")
  .catch(err => {console.log(err); "no"})
  return stat
}

/**
 * @param {string} pubid
 */
export async function deleteSource(pubid) {
  let stat = await fetch(get(api_url) + "/rss/" + pubid, { 
    "method": "DELETE",
  }).then(res => res.status === 200? "ok": "no")
  .catch(err => {console.log(err); "no"})
  return stat
}

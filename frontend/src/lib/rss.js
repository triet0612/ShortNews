import {api_url} from "./index.js"
import { get } from "svelte/store"

export class NewsSource {
  /**
   * @param {string} pubID 
   * @param {string} pub 
   * @param {string} link 
   * @param {string} lang 
   */
  constructor(pubID, pub, link, lang) {
    this.pubID = pubID
    this.pub = pub
    this.link = link
    this.lang = lang
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
      v["PublisherID"], v["Publisher"], v["Link"], v["Language"],
    ))
  })
  return options
}